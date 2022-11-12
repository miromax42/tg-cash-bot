package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/logtags"
	"github.com/go-testfixtures/testfixtures/v3"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	gruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"go.nhat.io/otelsql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency/fakeexchange"
	_ "gitlab.ozon.dev/miromaxxs/telegram-bot/doc/statik"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/gapi"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/pb"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo/cache/redis"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/sender"
	kafkasender "gitlab.ozon.dev/miromaxxs/telegram-bot/sender/kafka"
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
		bot          *tele.Bot
		db           *ent.Client
		cache        *redis.Cache
		exchange     currency.Exchange
		srv          *telegram.Server
		tp           *tracesdk.TracerProvider
		reportSender sender.ReportSender
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

	init.Go(func() (err error) {
		cache, err = redis.NewCache(initCtx, cfg.Cache)

		return errors.Wrap(err, "cache")
	})

	init.Go(func() (err error) {
		bot, err = createTelegramConnection(cfg.Telegram, log)

		return errors.Wrap(err, "telegram")
	})

	init.Go(func() (err error) {
		reportSender = kafkasender.NewWriter(cfg.Kafka)

		return nil
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

		srv, err = telegram.NewServer(workCtx, log, bot, expense, personalSettings, cache, exchange, reportSender)
		if err != nil {
			return err
		}

		log.Info(workCtx, "bot started")
		srv.Start()

		return nil
	})

	work.Go(func() error {
		return runGRPCServer(workCtx, cfg.GRPC, log, bot)
	})

	work.Go(func() error {
		return runGatewayServer(workCtx, cfg.HTTP, log, bot)
	})

	work.Go(func() error {
		return runMetricsServer(workCtx, cfg.HTTP, log)
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

	if err := work.Wait(); err != nil {
		log.Panicf(workCtx, "gracefull stop: %s \n", err)
	} else {
		log.Info(workCtx, "gracefully stopped!")
	}
}

func createTelegramConnection(cfg util.ConfigTelegram, log logger.Logger) (*tele.Bot, error) {
	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: time.Second},
		OnError: func(err error, c tele.Context) {
			log.Error(telegram.RequestContext(c), fmt.Sprintf("%+v", err))
		},
	}

	return tele.NewBot(pref)
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

func runGRPCServer(ctx context.Context, cfg util.ConfigGRPC, log logger.Logger, bot *tele.Bot) error {
	server, err := gapi.New(bot, log)
	if err != nil {
		return errors.Wrap(err, "create server")
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(grpc_validator.UnaryServerInterceptor()),
		),
	)

	pb.RegisterBotSendServer(grpcServer, server)
	reflection.Register(grpcServer) // explore server

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return errors.Wrap(err, "cannot create listener: %s")
	}

	go func() {
		<-ctx.Done()
		grpcServer.Stop()
	}()

	log.Printf(ctx, "stating GRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		return errors.Wrap(err, "cannot start GRPC server: %s")
	}

	return nil
}

func runGatewayServer(ctx context.Context, cfg util.ConfigHTTP, log logger.Logger, bot *tele.Bot) error {
	server, err := gapi.New(bot, log)
	if err != nil {
		return errors.Wrap(err, "cannot create server")
	}

	grpcMux := gruntime.NewServeMux(
		gruntime.WithMarshalerOption(gruntime.MIMEWildcard, &gruntime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	err = pb.RegisterBotSendHandlerServer(ctx, grpcMux, server)
	if err != nil {
		return errors.Wrap(err, "cannot register server handler")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		return errors.Wrap(err, "cannot create statik fs")
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return errors.Wrap(err, "cannot create listener")
	}

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	log.Printf(ctx, "stating HTTP gateway server on %s", listener.Addr().String())
	err = http.Serve(listener, mux) //nolint:gosec
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		return errors.Wrap(err, "cannot start http-gateway server")
	}

	return nil
}

func runMetricsServer(ctx context.Context, cfg util.ConfigHTTP, log logger.Logger) error {
	const metricsEndpoint = "/metrics"

	metrics := &http.Server{Addr: fmt.Sprintf(":%d", cfg.MetricsPort)} //nolint:gosec
	http.Handle(metricsEndpoint, promhttp.Handler())

	go func() {
		<-ctx.Done()
		metrics.Close()
	}()

	log.Printf(ctx, "stating HTTP Metrics server on [::]:%d%s", cfg.MetricsPort, metricsEndpoint)
	if err := metrics.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
