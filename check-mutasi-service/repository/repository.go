package repository

import (
	"context"
	"database/sql"

	"github.com/izzanzahrial/jojonomic-interview/entities"
)

type CheckMutation struct {
	db *sql.DB
}

func NewCheckMutation(db *sql.DB) *CheckMutation {
	return &CheckMutation{
		db: db,
	}
}

func (cm *CheckMutation) GetListTransactionByDate(ctx context.Context, norek string, startDate, endDate int32) ([]*entities.Transaction, error) {
	tx, err := cm.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		SELECT extract(epoch from timestamp with time zone created_at), 
		types, gram, topup_price, buyback_price, saldo
		FROM transactions
		WHERE norek = $1 AND created_at BETWEEN to_timestamp($2) AND to_timestamp($3) 
		ORDER BY created_at DESC`

	rows, err := tx.QueryContext(ctx, query, norek, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []*entities.Transaction{}
	for rows.Next() {
		var transaction entities.Transaction
		err := rows.Scan(&transaction.Date, &transaction.Type, &transaction.Gram, &transaction.TopupPrice,
			&transaction.BuybackPrice, &transaction.Saldo)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transactions, nil
}
