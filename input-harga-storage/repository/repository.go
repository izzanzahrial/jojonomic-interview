package repository

import (
	"context"
	"database/sql"

	"github.com/izzanzahrial/jojonomic-interview/entities"
)

type PriceStorage struct {
	db *sql.DB
}

func NewPriceStorage(db *sql.DB) *PriceStorage {
	return &PriceStorage{
		db: db,
	}
}

func (ps *PriceStorage) Insert(ctx context.Context, data *entities.Price) error {
	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO prices (admin_id, topup_price, buyback_price) VALUES ($1, $2, $3)`
	result, err := tx.ExecContext(ctx, query, data.AdminID, data.TopupPrice, data.BuybackPrice)
	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
