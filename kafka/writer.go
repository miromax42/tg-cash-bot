package kafka

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	kf "github.com/segmentio/kafka-go"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util/metrics"
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
	start := time.Now()

	err := errors.Wrapf(
		w.kr.WriteMessages(ctx, kf.Message{Value: value}),
		"kafka",
	)

	metrics.KafkaWriteDuration.Observe(
		time.Since(start).Seconds(),
	)
	metrics.KafkaWriteCount.Inc()

	return err
}
