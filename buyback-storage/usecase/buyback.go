package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/izzanzahrial/jojonomic-interview/buyback-storage/repository"
	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/entities"
)

type BuybackStorage struct {
	consumer   common.KafkaConsumer
	repository *repository.BuybackStorage
}

func NewBuybackStorage(consumer common.KafkaConsumer, repo *repository.BuybackStorage) *BuybackStorage {
	return &BuybackStorage{
		consumer:   consumer,
		repository: repo,
	}
}

func (bs *BuybackStorage) Consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context canceled, stopping consumer")
			return
		default:
			msg, err := bs.consumer.ReadMessage(ctx)
			if err != nil {
				fmt.Println("Error when consuming message: ", err)
				continue
			}
			fmt.Println("Recieved message: ", string(msg.Value))

			var buyback entities.Buyback
			if err := json.Unmarshal(msg.Value, &buyback); err != nil {
				fmt.Println("Error when decoding json: ", err)
				continue
			}

			fmt.Println("Recieved message: ", buyback)
			err = bs.repository.InsertAndUpdate(ctx, &buyback)
			if err != nil {
				fmt.Println("Error when insert and update data to the database: ", err)
				continue
			}
		}
	}
}
