package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/izzanzahrial/jojonomic-interview/common"
)

type CheckSaldo struct {
	db *sql.DB
}

func NewCheckSaldo(db *sql.DB) *CheckSaldo {
	return &CheckSaldo{
		db: db,
	}
}

func (cs *CheckSaldo) GetSaldo(ctx context.Context, norek string) (float32, error) {
	tx, err := cs.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var saldo float32

	query := `SELECT saldo FROM accounts WHERE norek = $1`
	err = tx.QueryRowContext(ctx, query, norek).Scan(&saldo)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, common.ErrRecordNotFound
		default:
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return saldo, nil
}
