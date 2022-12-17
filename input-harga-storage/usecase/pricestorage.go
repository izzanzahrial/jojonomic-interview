package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/entities"
	"github.com/izzanzahrial/jojonomic-interview/input-harga-storage/repository"
)

type PriceStorage struct {
	consumer   common.KafkaConsumer
	repository *repository.PriceStorage
}

func NewPriceStorage(consumer common.KafkaConsumer, repo *repository.PriceStorage) *PriceStorage {
	return &PriceStorage{
		consumer:   consumer,
		repository: repo,
	}
}

func (p *PriceStorage) Insert(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context canceled, stopping consumer")
			return
		default:
			msg, err := p.consumer.ReadMessage(ctx)
			if err != nil {
				fmt.Println("Error when consuming message: ", err)
				continue
			}
			fmt.Println("Recieved message: ", string(msg.Value))

			var price entities.Price
			if err := json.Unmarshal(msg.Value, &price); err != nil {
				fmt.Println("Error when decoding json: ", err)
				continue
			}

			fmt.Println("Recieved message: ", price)
			err = p.repository.Insert(ctx, &price)
			if err != nil {
				fmt.Println("Error when insert price data to the database: ", err)
				continue
			}
		}
	}
}
