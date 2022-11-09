package kafkasender

import (
	"context"

	"github.com/cockroachdb/errors"
	"google.golang.org/protobuf/proto"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/kafka"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/pb"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

const (
	topicReportRequest = "route256.expenses-bot.report-request"
)

type Writer struct {
	kr *kafka.Writer
}

func NewWriter(cfg util.ConfigKafka) *Writer {
	writer := kafka.NewWriter(topicReportRequest, cfg)

	return &Writer{writer}
}

func (w *Writer) ReportSend(ctx context.Context, req *pb.ReportRequest) error {
	bytes, err := proto.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "marshal")
	}

	return errors.Wrap(
		w.kr.SendBytes(ctx, bytes),
		"send",
	)
}
