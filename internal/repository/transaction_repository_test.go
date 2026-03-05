package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/transactions-platform/internal/models"
)

func TestTransactionRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewTransactionRepository(db)
	ctx := context.Background()

	tests := []struct {
		name          string
		accountID     string
		operationType models.OperationType
		amount        decimal.Decimal
		mockSetup     func()
		wantErr       bool
	}{
		{
			name:          "successful purchase transaction creation",
			accountID:     "550e8400-e29b-41d4-a716-446655440000",
			operationType: models.NormalPurchase,
			amount:        decimal.NewFromFloat(-123.45),
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "account_id", "amount", "event_date", "operation_type", "created_at"}).
					AddRow("tx-id", "550e8400-e29b-41d4-a716-446655440000", decimal.NewFromFloat(-123.45), time.Now(), "NORMAL_PURCHASE", time.Now())
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs("550e8400-e29b-41d4-a716-446655440000", models.NormalPurchase, decimal.NewFromFloat(-123.45)).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:          "successful credit voucher creation",
			accountID:     "550e8400-e29b-41d4-a716-446655440000",
			operationType: models.CreditVoucher,
			amount:        decimal.NewFromFloat(100.00),
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "account_id", "amount", "event_date", "operation_type", "created_at"}).
					AddRow("tx-id", "550e8400-e29b-41d4-a716-446655440000", decimal.NewFromFloat(100.00), time.Now(), "CREDIT_VOUCHER", time.Now())
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs("550e8400-e29b-41d4-a716-446655440000", models.CreditVoucher, decimal.NewFromFloat(100.00)).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:          "successful withdrawal creation",
			accountID:     "550e8400-e29b-41d4-a716-446655440000",
			operationType: models.Withdrawal,
			amount:        decimal.NewFromFloat(-50.00),
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "account_id", "amount", "event_date", "operation_type", "created_at"}).
					AddRow("tx-id", "550e8400-e29b-41d4-a716-446655440000", decimal.NewFromFloat(-50.00), time.Now(), "WITHDRAWAL", time.Now())
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs("550e8400-e29b-41d4-a716-446655440000", models.Withdrawal, decimal.NewFromFloat(-50.00)).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:          "successful installment purchase creation",
			accountID:     "550e8400-e29b-41d4-a716-446655440000",
			operationType: models.PurchaseWithInstallments,
			amount:        decimal.NewFromFloat(-300.00),
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "account_id", "amount", "event_date", "operation_type", "created_at"}).
					AddRow("tx-id", "550e8400-e29b-41d4-a716-446655440000", decimal.NewFromFloat(-300.00), time.Now(), "PURCHASE_WITH_INSTALLMENTS", time.Now())
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs("550e8400-e29b-41d4-a716-446655440000", models.PurchaseWithInstallments, decimal.NewFromFloat(-300.00)).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:          "database error",
			accountID:     "550e8400-e29b-41d4-a716-446655440000",
			operationType: models.NormalPurchase,
			amount:        decimal.NewFromFloat(-123.45),
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs("550e8400-e29b-41d4-a716-446655440000", models.NormalPurchase, decimal.NewFromFloat(-123.45)).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name:          "foreign key violation - account not found",
			accountID:     "non-existent",
			operationType: models.NormalPurchase,
			amount:        decimal.NewFromFloat(-123.45),
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs("non-existent", models.NormalPurchase, decimal.NewFromFloat(-123.45)).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			transaction, err := repo.Create(ctx, tt.accountID, tt.operationType, tt.amount)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, transaction)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, transaction)
				assert.Equal(t, tt.accountID, transaction.AccountID)
				assert.Equal(t, tt.operationType, transaction.OperationType)
				assert.True(t, tt.amount.Equal(transaction.Amount), "expected amount %s but got %s", tt.amount.String(), transaction.Amount.String())
				assert.NotEmpty(t, transaction.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
