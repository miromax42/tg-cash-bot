package main

import (
	"context"
	"net/http"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/logtags"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.nhat.io/otelsql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency/fakeexchange"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util/logger"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo/database"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func main() { //nolint:funlen
	var (
		db       *ent.Client
		exchange currency.Exchange
		srv      *telegram.Server
		tp       *tracesdk.TracerProvider
		metricts *http.Server
	)

	ctx := logtags.AddTag(context.Background(), "golang.version", runtime.Version())

	mainCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := logger.New()

	cfg, err := util.NewConfig()
	if err != nil {
		log.Panic(mainCtx, err)
	}

	if tp, err = tracerProvider(cfg.Tracing); err != nil {
		log.Panic(mainCtx, errors.Wrap(err, "init trace provider"))
	}
	otel.SetTracerProvider(tp)

	init, initCtx := errgroup.WithContext(mainCtx)

	init.Go(func() (err error) {
		driver, err := openDB(cfg.DB.URL)
		if err != nil {
			return errors.Wrap(err, "open db")
		}

		db, err = migrateDB(initCtx, driver)
		if err != nil {
			return errors.Wrap(err, "ent migration")
		}

		err = initFixtures(driver, cfg.DB.TestUserID)
		if err != nil {
			return errors.Wrap(err, "db fixtures")
		}

		return nil
	})

	init.Go(func() (err error) {
		exchange = fakeexchange.Exchange{}

		return errors.Wrap(err, "exchange")
	})

	if err = init.Wait(); err != nil {
		log.Panicf(initCtx, errors.Wrap(err, "init").Error())
	}

	work, workCtx := errgroup.WithContext(mainCtx)

	work.Go(func() (err error) {
		db.Use(func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				ctx, span := otel.Tracer(util.RequestTrace).Start(ctx, "mutation",
					trace.WithAttributes(
						attribute.Stringer("op", m.Op()),
						attribute.String("type", m.Type()),
					))
				defer span.End()

				return next.Mutate(ctx, m)
			})
		})

		expense := database.NewExpense(db)
		personalSettings := database.NewPersonalSettings(db)

		srv, err = telegram.NewServer(workCtx, cfg.Telegram, log, expense, personalSettings, exchange)
		if err != nil {
			return err
		}

		log.Info(workCtx, "bot started")
		srv.Start()

		return nil
	})

	work.Go(func() error {
		metricts = &http.Server{Addr: ":2112"} //nolint:gosec
		http.Handle("/metrics", promhttp.Handler())

		if !errors.Is(metricts.ListenAndServe(), http.ErrServerClosed) {
			return err
		}

		return nil
	})

	work.Go(func() error {
		<-workCtx.Done()

		srv.Stop()

		return errors.Wrap(db.Close(), "close db")
	})

	work.Go(func() error {
		<-workCtx.Done()

		tpCtx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		return tp.Shutdown(tpCtx)
	})

	work.Go(func() error {
		<-workCtx.Done()

		sCtx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		return metricts.Shutdown(sCtx)
	})

	if err := work.Wait(); err != nil {
		log.Panicf(workCtx, "gracefull stop: %s \n", err)
	} else {
		log.Info(workCtx, "gracefully stopped!")
	}
}

func migrateDB(ctx context.Context, driver *sql.Driver) (*ent.Client, error) {
	db := ent.NewClient(ent.Driver(driver))

	if err := db.Schema.Create(ctx); err != nil {
		return nil, errors.Wrap(err, "migrate")
	}

	return db, nil
}

func initFixtures(driver *sql.Driver, id int64) error {
	if id == 0 {
		return nil
	}

	fixtures, err := testfixtures.New(
		testfixtures.Template(),
		testfixtures.TemplateData(map[string]interface{}{
			"ID": id,
		}),
		testfixtures.Database(driver.DB()),
		testfixtures.Dialect(driver.Dialect()),
		testfixtures.FilesMultiTables("ent/fixtures/test_user_seed.yml"),
	)
	if err != nil {
		return err
	}

	return fixtures.Load()
}

func tracerProvider(cfg util.ConfigTracing) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.URL)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("telegram bot"),
		)),
	)

	return tp, nil
}

func openDB(dsn string) (*sql.Driver, error) {
	driverName, err := otelsql.Register("postgres",
		otelsql.TraceQueryWithoutArgs(),
		otelsql.TraceRowsClose(),
		otelsql.TraceRowsAffected(),
		otelsql.WithSystem(semconv.DBSystemPostgreSQL),
	)
	if err != nil {
		return nil, err
	}

	return sql.Open(driverName, dsn)
}
