package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/logtags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	_ "github.com/lib/pq"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/kafka"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/pb"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo/database"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util/logger"
)

const (
	topicReportRequest = "route256.expenses-bot.report-request"
	groupReportRequest = "reportGroup1"
)

//nolint:funlen
func main() {
	ctx := context.Background()

	log := logger.New()

	cfg, err := util.NewConfig()
	if err != nil {
		log.Panic(ctx, err)
	}

	// init db
	log.Printf(ctx, "connecting to BD:%q", cfg.DB.URL)

	dbClient, err := ent.Open("postgres", cfg.DB.URL)
	if err != nil {
		log.Panic(ctx, err)
	}

	expenses := database.NewExpense(dbClient)

	// init grpc
	log.Printf(ctx, "connecting to GRPC-server:%q", cfg.GRPC.Address)

	conn, err := grpc.Dial(cfg.GRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panic(ctx, err)
	}
	defer conn.Close()

	client := pb.NewBotSendClient(conn)

	// kafka
	log.Printf(ctx, "connecting to Kafka:%q", cfg.Kafka.Address)

	reader := kafka.NewReader(topicReportRequest, groupReportRequest, cfg.Kafka)
	for msg := range reader.Read(ctx) {
		if msg.Err != nil {
			log.Error(ctx, err)

			continue
		}

		var req pb.ReportRequest
		if err = proto.Unmarshal(msg.Value, &req); err != nil {
			log.Error(ctx, err)

			continue
		}

		listUserExpense, err := expenses.ListUserExpense(ctx, repo.ListUserExpenseReq{
			UserID:   req.GetUserId(),
			FromTime: req.GetStartTime().AsTime(),
			ToTime:   req.GetEndTime().AsTime(),
		})
		if err != nil {
			log.Error(ctx, err)

			continue
		}

		rsp, err := client.SendReport(ctx, &pb.SendMessageRequest{
			UserId:  req.UserId,
			Message: renderReport(listUserExpense, req.GetMultiplier()),
		})
		if err != nil {
			log.Error(ctx, err)
		}

		mCtx := logtags.AddTag(ctx, "service", "renderer")
		if rsp.Success {
			log.Info(mCtx, "success")
		} else {
			log.Warn(mCtx, "no success")
		}
	}
}

func renderReport(listExp repo.ListUserExpenseResp, multiplier float64) string {
	b := strings.Builder{}

	if len(listExp) == 0 {
		return "You have no expenses!"
	}

	b.WriteString("Your expenses:\n")
	for i := range listExp {
		line := fmt.Sprintf("%d. %s - %0.2f\n", i+1, listExp[i].Category, listExp[i].Sum*multiplier)
		b.WriteString(line)
	}

	return b.String()
}
