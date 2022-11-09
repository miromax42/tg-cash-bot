package kafka

import (
	"context"

	kf "github.com/segmentio/kafka-go"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type Writer struct {
	kr *kf.Writer
}

func NewWriter(topic string, cfg util.ConfigKafka) *Writer {
	writer := &kf.Writer{
		Addr:  kf.TCP(cfg.Address),
		Topic: topic,
	}

	return &Writer{writer}
}

func (w *Writer) SendBytes(ctx context.Context, value []byte) error {
	return w.kr.WriteMessages(ctx, kf.Message{Value: value})
}
