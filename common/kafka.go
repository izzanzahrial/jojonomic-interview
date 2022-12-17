package common

import (
	"github.com/izzanzahrial/jojonomic-interview/config"
	"github.com/segmentio/kafka-go"
)

func NewKafkaWriter(cfg config.Kafka) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(cfg.Hosts...),
		Topic: cfg.Topic,
		Balancer: &kafka.CRC32Balancer{
			Consistent: true,
		},
		RequiredAcks: kafka.RequireAll,
	}
}

func NewKafkaConsumer(cfg config.Kafka) *kafka.Reader {
	return kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:     cfg.Hosts,
			Topic:       cfg.Topic,
			GroupID:     "jojonomic-price-main-consumer",
			StartOffset: kafka.LastOffset,
		},
	)
}
