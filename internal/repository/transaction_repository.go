package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/transactions-platform/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, accountID string, operationType models.OperationType, amount decimal.Decimal) (*models.Transaction, error) {
	transaction := &models.Transaction{}

	query := `
		INSERT INTO transactions (account_id, operation_type, amount)
		VALUES ($1, $2, $3)
		RETURNING id, account_id, amount, event_date, operation_type, created_at
	`

	err := r.db.QueryRowContext(ctx, query, accountID, operationType, amount).Scan(
		&transaction.ID,
		&transaction.AccountID,
		&transaction.Amount,
		&transaction.EventDate,
		&transaction.OperationType,
		&transaction.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

