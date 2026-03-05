package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/transactions-platform/internal/models"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, documentNumber string) (*models.Account, error) {
	account := &models.Account{}

	query := `
		INSERT INTO accounts (document_number)
		VALUES ($1)
		RETURNING id, document_number, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query, documentNumber).Scan(
		&account.ID,
		&account.DocumentNumber,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id string) (*models.Account, error) {
	account := &models.Account{}

	query := `
		SELECT id, document_number, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.DocumentNumber,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

func (r *AccountRepository) GetByDocumentNumber(ctx context.Context, documentNumber string) (*models.Account, error) {
	account := &models.Account{}

	query := `
		SELECT id, document_number, created_at, updated_at
		FROM accounts
		WHERE document_number = $1
	`

	err := r.db.QueryRowContext(ctx, query, documentNumber).Scan(
		&account.ID,
		&account.DocumentNumber,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found, return nil without error
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}
