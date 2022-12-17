package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/entities"
	"github.com/segmentio/kafka-go"
)

type Topup struct {
	publisher common.KafkaPublisher
}

func NewTopup(publisher common.KafkaPublisher) *Topup {
	return &Topup{
		publisher: publisher,
	}
}

func (t *Topup) Publish(ctx context.Context, reffID string, topup *entities.Topup) error {
	val, err := json.Marshal(topup)
	if err != nil {
		return err
	}
	fmt.Println(string(val))

	message := kafka.Message{
		Key:   []byte(reffID),
		Value: val,
	}
	fmt.Println(message)

	err = t.publisher.WriteMessages(ctx, message)
	if err != nil {
		return err
	}
	fmt.Println("Send message:", string(message.Value))

	return nil
}
