package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/izzanzahrial/jojonomic-interview/buyback-service/handler"
	"github.com/izzanzahrial/jojonomic-interview/buyback-service/usecase"
	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/config"
	"github.com/joho/godotenv"
	"github.com/teris-io/shortid"
)

func main() {
	if err := godotenv.Load("../../sample.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	var config config.Kafka
	config.Hosts = append(config.Hosts, os.Getenv("BUYBACK_KAFKA_HOSTS"))
	config.Topic = os.Getenv("BUYBACK_KAFKA_TOPIC")
	serviceID := os.Getenv("BUYBACK_SERVICE_ID")
	servicePort := os.Getenv("BUYBACK_SERVICE_PORT")

	writer := common.NewKafkaWriter(config)
	defer writer.Close()

	service := usecase.NewBuyback(writer)

	numID, err := strconv.Atoi(serviceID)
	if err != nil {
		log.Fatal("Error convert service ID to int")
	}

	sid, err := shortid.New(uint8(numID), shortid.DefaultABC, 2342)
	if err != nil {
		log.Fatal("Error failed construct new instance of short ID")
	}

	handler := handler.NewBuyback(service, sid)

	r := mux.NewRouter()
	r.HandleFunc("/api/buyback", handler.Publish).Methods("POST")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	fmt.Println("Server is up and running")

	log.Fatal(http.ListenAndServe(":"+servicePort, r))

	<-sig

	os.Exit(0)
}
