package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/entities"
)

type BuybackStorage struct {
	db *sql.DB
}

func NewBuybackStorage(db *sql.DB) *BuybackStorage {
	return &BuybackStorage{
		db: db,
	}
}

func (bs *BuybackStorage) InsertAndUpdate(ctx context.Context, data *entities.Buyback) error {
	tx, err := bs.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var saldo float32

	query := `SELECT saldo FROM accounts WHERE norek = $1`
	err = tx.QueryRowContext(ctx, query, data.Norek).Scan(&saldo)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return common.ErrRecordNotFound
		default:
			return err
		}
	}

	query = `
		INSERT INTO transactions (norek, types, gram, topup_price, buyback_price, saldo) 
		VALUES ($1, $2, $3, $4, $5, $6)`
	result, err := tx.ExecContext(ctx, query, data.Norek, data.Gram, data.Topup, data.Harga, saldo)
	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return err
	}

	query = `UPDATE accounts SET saldo = saldo - $1, updated_at = NOW() WHERE norek = $2`
	result, err = tx.ExecContext(ctx, query, data.Gram, data.Norek)
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
