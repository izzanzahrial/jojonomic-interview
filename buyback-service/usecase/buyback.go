package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/entities"
	"github.com/segmentio/kafka-go"
)

type Buyback struct {
	publisher common.KafkaPublisher
}

func NewBuyback(publisher common.KafkaPublisher) *Buyback {
	return &Buyback{
		publisher: publisher,
	}
}

func (b *Buyback) Publish(ctx context.Context, reffID string, buyback *entities.Buyback) error {
	val, err := json.Marshal(buyback)
	if err != nil {
		return err
	}
	fmt.Println(string(val))

	message := kafka.Message{
		Key:   []byte(reffID),
		Value: val,
	}
	fmt.Println(message)

	err = b.publisher.WriteMessages(ctx, message)
	if err != nil {
		return err
	}
	fmt.Println("Send message:", string(message.Value))

	return nil
}
