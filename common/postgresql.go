package common

import (
	"database/sql"
	"time"

	"github.com/izzanzahrial/jojonomic-interview/config"
)

func NewPostgreSQL(cfg config.PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	// This will result EOF in accountntransaction database
	// even though the error still there when we try to call the database
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// err = db.PingContext(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println("after")

	return db, nil
}
