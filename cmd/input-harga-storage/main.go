package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/config"
	"github.com/izzanzahrial/jojonomic-interview/input-harga-storage/repository"
	"github.com/izzanzahrial/jojonomic-interview/input-harga-storage/usecase"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load("../../sample.env"); err != nil {
		fmt.Println("Error loading .env file")
	}

	var postgresCfg config.PostgresConfig
	postgresCfg.DSN = os.Getenv("PRICE_POSTGRES_DSN")
	postgresCfg.MaxIdleTime = os.Getenv("PRICE_POSTGRES_MAX_IDLE_TIME")

	maxOpenConns := os.Getenv("PRICE_POSTGRES_MAX_OPEN_CONNS")
	postgresCfg.MaxOpenConns, _ = strconv.Atoi(maxOpenConns)

	maxIdleConns := os.Getenv("PRICE_POSTGRES_MAX_IDLE_CONNS")
	postgresCfg.MaxIdleConns, _ = strconv.Atoi(maxIdleConns)

	db, err := common.NewPostgreSQL(postgresCfg)
	if err != nil {
		fmt.Println("Error failed to established connection with the database", err)
	}

	var kafkaCfg config.Kafka
	kafkaCfg.Hosts = append(kafkaCfg.Hosts, os.Getenv("PRICE_KAFKA_HOSTS"))
	kafkaCfg.Topic = os.Getenv("PRICE_KAFKA_TOPIC")

	repository := repository.NewPriceStorage(db)
	consumer := common.NewKafkaConsumer(kafkaCfg)
	pricestorage := usecase.NewPriceStorage(consumer, repository)
	defer consumer.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		pricestorage.Insert(ctx)
	}()

	fmt.Println("Consumer is up and running")
	<-sig
	cancel()

	os.Exit(0)
}
