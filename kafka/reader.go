package kafka

import (
	"context"

	"github.com/cockroachdb/errors"
	kf "github.com/segmentio/kafka-go"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type Reader struct {
	kr *kf.Reader
}

type Message struct {
	kf.Message
	Err error
}

func NewReader(topic, group string, cfg util.ConfigKafka) *Reader {
	reader := kf.NewReader(
		kf.ReaderConfig{
			Brokers: []string{cfg.Address},
			Topic:   topic,
			GroupID: group,
		})

	return &Reader{reader}
}

func (r *Reader) Read(ctx context.Context) <-chan Message {
	return r.commit(ctx, r.fetch(ctx))
}

func (r *Reader) fetch(ctx context.Context) <-chan Message {
	msgStream := make(chan Message)

	go func() {
		defer close(msgStream)

		for {
			kfMsg, err := r.kr.FetchMessage(ctx)

			msg := Message{
				Message: kfMsg,
				Err:     err,
			}

			select {
			case <-ctx.Done():
				return
			case msgStream <- msg:
			}
		}
	}()

	return msgStream
}

func (r *Reader) commit(ctx context.Context, inputStream <-chan Message) <-chan Message {
	committedStream := make(chan Message)

	go func() {
		defer close(committedStream)

		for {
			for v := range util.OrDone(ctx.Done(), inputStream) {
				err := r.kr.CommitMessages(ctx, v.Message)
				v.Err = errors.CombineErrors(v.Err, err)

				select {
				case <-ctx.Done():
					return
				case committedStream <- v:
				}
			}
		}

	}()

	return committedStream
}
