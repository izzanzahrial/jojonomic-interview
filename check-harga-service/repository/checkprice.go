package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/izzanzahrial/jojonomic-interview/common"
	"github.com/izzanzahrial/jojonomic-interview/entities"
)

type CheckPrice struct {
	db *sql.DB
}

func NewCheckPrice(db *sql.DB) *CheckPrice {
	return &CheckPrice{
		db: db,
	}
}

func (cp *CheckPrice) CheckGoldPrice(ctx context.Context) (*entities.Price, error) {
	tx, err := cp.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `SELECT topup_price, buyback_price FROM prices ORDER BY created_at DESC LIMIT 1`

	var price entities.Price

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := cp.db.QueryRowContext(ctx, query).Scan(&price.TopupPrice, &price.BuybackPrice); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, common.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &price, nil
}
