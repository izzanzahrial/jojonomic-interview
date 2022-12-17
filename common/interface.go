package common

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaPublisher interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type KafkaConsumer interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
}
