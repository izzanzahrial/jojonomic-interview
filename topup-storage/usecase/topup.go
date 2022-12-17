package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/entities"
	"github.com/izzanzahrial/jojonomic-interview/topup-storage/repository"
)

type TopupStorage struct {
	consumer   common.KafkaConsumer
	repository *repository.TopupStorage
}

func NewTopupStorage(consumer common.KafkaConsumer, repo *repository.TopupStorage) *TopupStorage {
	return &TopupStorage{
		consumer:   consumer,
		repository: repo,
	}
}

func (ts *TopupStorage) Consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context canceled, stopping consumer")
			return
		default:
			msg, err := ts.consumer.ReadMessage(ctx)
			if err != nil {
				fmt.Println("Error when consuming message: ", err)
				continue
			}
			fmt.Println("Recieved message: ", string(msg.Value))

			var topup entities.Topup
			if err := json.Unmarshal(msg.Value, &topup); err != nil {
				fmt.Println("Error when decoding json: ", err)
				continue
			}

			fmt.Println("Recieved message: ", topup)
			err = ts.repository.InsertAndUpdate(ctx, &topup)
			if err != nil {
				fmt.Println("Error when insert and update data to the database: ", err)
				continue
			}
		}
	}
}
