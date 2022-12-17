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
	"github.com/izzanzahrial/jojonomic-interview/check-saldo-service/handler"
	"github.com/izzanzahrial/jojonomic-interview/check-saldo-service/repository"
	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/config"
	"github.com/joho/godotenv"
	"github.com/teris-io/shortid"

	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load("../../sample.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	var postgresCfg config.PostgresConfig
	postgresCfg.DSN = os.Getenv("ACCOUNT_TRANSACTION_POSTGRES_DSN")
	postgresCfg.MaxIdleTime = os.Getenv("ACCOUNT_TRANSACTION_POSTGRES_MAX_IDLE_TIME")

	maxOpenConns := os.Getenv("ACCOUNT_TRANSACTION_POSTGRES_MAX_OPEN_CONNS")
	postgresCfg.MaxOpenConns, _ = strconv.Atoi(maxOpenConns)

	maxIdleConns := os.Getenv("ACCOUNT_TRANSACTION_POSTGRES_MAX_IDLE_CONNS")
	postgresCfg.MaxIdleConns, _ = strconv.Atoi(maxIdleConns)

	serviceID := os.Getenv("SALDO_SERVICE_ID")
	servicePort := os.Getenv("SALDO_SERVICE_PORT")

	db, err := common.NewPostgreSQL(postgresCfg)
	if err != nil {
		fmt.Println("Error failed to established connection with the database", err)
	}

	numID, err := strconv.Atoi(serviceID)
	if err != nil {
		log.Fatal("Error convert service ID to int")
	}

	sid, err := shortid.New(uint8(numID), shortid.DefaultABC, 2342)
	if err != nil {
		log.Fatal("Error failed construct new instance of short ID")
	}

	repository := repository.NewCheckSaldo(db)
	handler := handler.NewCheckSaldo(repository, sid)

	r := mux.NewRouter()
	r.HandleFunc("/api/saldo", handler.GetSaldo).Methods("GET")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	fmt.Println("Server is up and running")

	log.Fatal(http.ListenAndServe(":"+servicePort, r))

	<-sig

	os.Exit(0)
}
