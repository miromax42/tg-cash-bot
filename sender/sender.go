package sender

import (
	"context"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/pb"
)

type ReportSender interface {
	ReportSend(ctx context.Context, req *pb.ReportRequest) error
}
