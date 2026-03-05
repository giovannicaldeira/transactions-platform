package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/transactions-platform/internal/repository"
)

func TestAccountService_CreateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewAccountRepository(db)
	service := NewAccountService(repo)
	ctx := context.Background()

	tests := []struct {
		name           string
		documentNumber string
		mockSetup      func()
		wantErr        bool
		errContains    string
	}{
		{
			name:           "successful creation",
			documentNumber: "12345678900",
			mockSetup: func() {
				// Check for existing - not found
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE document_number").
					WithArgs("12345678900").
					WillReturnError(sql.ErrNoRows)

				// Create new account
				rows := sqlmock.NewRows([]string{"id", "document_number", "created_at", "updated_at"}).
					AddRow("550e8400-e29b-41d4-a716-446655440000", "12345678900", time.Now(), time.Now())
				mock.ExpectQuery("INSERT INTO accounts").
					WithArgs("12345678900").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:           "duplicate document number",
			documentNumber: "12345678900",
			mockSetup: func() {
				// Return existing account
				rows := sqlmock.NewRows([]string{"id", "document_number", "created_at", "updated_at"}).
					AddRow("550e8400-e29b-41d4-a716-446655440000", "12345678900", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE document_number").
					WithArgs("12345678900").
					WillReturnRows(rows)
			},
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name:           "database error on check",
			documentNumber: "12345678900",
			mockSetup: func() {
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE document_number").
					WithArgs("12345678900").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr:     true,
			errContains: "failed to check existing account",
		},
		{
			name:           "database error on create",
			documentNumber: "12345678900",
			mockSetup: func() {
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE document_number").
					WithArgs("12345678900").
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery("INSERT INTO accounts").
					WithArgs("12345678900").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr:     true,
			errContains: "failed to create account",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			account, err := service.CreateAccount(ctx, tt.documentNumber)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, account)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
				assert.Equal(t, tt.documentNumber, account.DocumentNumber)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAccountService_GetAccountByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewAccountRepository(db)
	service := NewAccountService(repo)
	ctx := context.Background()

	tests := []struct {
		name        string
		accountID   string
		mockSetup   func()
		wantErr     bool
		errContains string
	}{
		{
			name:      "account found",
			accountID: "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "document_number", "created_at", "updated_at"}).
					AddRow("550e8400-e29b-41d4-a716-446655440000", "12345678900", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:      "account not found",
			accountID: "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func() {
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr:     true,
			errContains: "account not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			account, err := service.GetAccountByID(ctx, tt.accountID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
