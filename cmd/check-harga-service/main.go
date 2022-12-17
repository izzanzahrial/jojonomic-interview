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
	"github.com/izzanzahrial/jojonomic-interview/check-harga-service/handler"
	"github.com/izzanzahrial/jojonomic-interview/check-harga-service/repository"
	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/config"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load("../../sample.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	var postgresCfg config.PostgresConfig
	postgresCfg.DSN = os.Getenv("PRICE_POSTGRES_DSN")
	postgresCfg.MaxIdleTime = os.Getenv("PRICE_POSTGRES_MAX_IDLE_TIME")

	maxOpenConns := os.Getenv("PRICE_POSTGRES_MAX_OPEN_CONNS")
	postgresCfg.MaxOpenConns, _ = strconv.Atoi(maxOpenConns)

	maxIdleConns := os.Getenv("PRICE_POSTGRES_MAX_IDLE_CONNS")
	postgresCfg.MaxIdleConns, _ = strconv.Atoi(maxIdleConns)

	servicePort := os.Getenv("CHECK_PRICE_SERVICE_PORT")

	db, err := common.NewPostgreSQL(postgresCfg)
	if err != nil {
		fmt.Println("Error failed to established connection with the database", err)
	}

	repository := repository.NewCheckPrice(db)
	handler := handler.NewCheckPrice(repository)

	r := mux.NewRouter()
	r.HandleFunc("/api/check-harga", handler.CheckGoldPrice).Methods("GET")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	fmt.Println("Server is up and running")

	log.Fatal(http.ListenAndServe(":"+servicePort, r))

	<-sig

	os.Exit(0)
}
