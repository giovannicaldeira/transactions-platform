package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewAccountRepository(db)
	ctx := context.Background()

	tests := []struct {
		name           string
		documentNumber string
		mockSetup      func()
		wantErr        bool
	}{
		{
			name:           "successful creation",
			documentNumber: "12345678900",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "document_number", "created_at", "updated_at"}).
					AddRow("550e8400-e29b-41d4-a716-446655440000", "12345678900", time.Now(), time.Now())
				mock.ExpectQuery("INSERT INTO accounts").
					WithArgs("12345678900").
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:           "database error",
			documentNumber: "12345678900",
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO accounts").
					WithArgs("12345678900").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			account, err := repo.Create(ctx, tt.documentNumber)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, account)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
				assert.Equal(t, tt.documentNumber, account.DocumentNumber)
				assert.NotEmpty(t, account.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAccountRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewAccountRepository(db)
	ctx := context.Background()

	tests := []struct {
		name      string
		accountID string
		mockSetup func()
		wantErr   bool
		errMsg    string
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
			wantErr: true,
			errMsg:  "account not found",
		},
		{
			name:      "database error",
			accountID: "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func() {
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			account, err := repo.GetByID(ctx, tt.accountID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
				assert.Equal(t, tt.accountID, account.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAccountRepository_GetByDocumentNumber(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewAccountRepository(db)
	ctx := context.Background()

	tests := []struct {
		name           string
		documentNumber string
		mockSetup      func()
		wantNil        bool
		wantErr        bool
	}{
		{
			name:           "account found",
			documentNumber: "12345678900",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "document_number", "created_at", "updated_at"}).
					AddRow("550e8400-e29b-41d4-a716-446655440000", "12345678900", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE document_number").
					WithArgs("12345678900").
					WillReturnRows(rows)
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name:           "account not found",
			documentNumber: "99999999999",
			mockSetup: func() {
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE document_number").
					WithArgs("99999999999").
					WillReturnError(sql.ErrNoRows)
			},
			wantNil: true,
			wantErr: false,
		},
		{
			name:           "database error",
			documentNumber: "12345678900",
			mockSetup: func() {
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE document_number").
					WithArgs("12345678900").
					WillReturnError(sql.ErrConnDone)
			},
			wantNil: true,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			account, err := repo.GetByDocumentNumber(ctx, tt.documentNumber)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.wantNil {
				assert.Nil(t, account)
			} else {
				assert.NotNil(t, account)
				assert.Equal(t, tt.documentNumber, account.DocumentNumber)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
