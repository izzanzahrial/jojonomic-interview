package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/entities"
	"github.com/segmentio/kafka-go"
)

type Inputprice struct {
	publisher common.KafkaPublisher
}

func NewInputPrice(publisher common.KafkaPublisher) *Inputprice {
	return &Inputprice{
		publisher: publisher,
	}
}

func (i *Inputprice) Publish(ctx context.Context, reffID string, price *entities.Price) error {
	val, err := json.Marshal(price)
	if err != nil {
		return err
	}
	fmt.Println(string(val))

	message := kafka.Message{
		Key:   []byte(reffID),
		Value: val,
	}
	fmt.Println(message)

	err = i.publisher.WriteMessages(ctx, message)
	if err != nil {
		return err
	}
	fmt.Println("Send message:", string(message.Value))

	return nil
}
