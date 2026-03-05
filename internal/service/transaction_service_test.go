package service

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
	"github.com/transactions-platform/internal/repository"
)

func TestTransactionService_CreateTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	service := NewTransactionService(transactionRepo, accountRepo)
	ctx := context.Background()

	tests := []struct {
		name          string
		accountID     string
		operationType models.OperationType
		amount        decimal.Decimal
		mockSetup     func()
		wantErr       bool
		errContains   string
		expectedSign  int // Expected sign of stored amount: -1, 0, or 1
	}{
		{
			name:          "successful purchase transaction",
			accountID:     "550e8400-e29b-41d4-a716-446655440000",
			operationType: models.NormalPurchase,
			amount:        decimal.NewFromFloat(123.45),
			mockSetup: func() {
				// Account exists check
				rows := sqlmock.NewRows([]string{"id", "document_number", "created_at", "updated_at"}).
					AddRow("550e8400-e29b-41d4-a716-446655440000", "12345678900", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnRows(rows)

				// Create transaction (amount should be negative)
				txRows := sqlmock.NewRows([]string{"id", "account_id", "amount", "event_date", "operation_type", "created_at"}).
					AddRow("tx-id", "550e8400-e29b-41d4-a716-446655440000", decimal.NewFromFloat(-123.45), time.Now(), "NORMAL_PURCHASE", time.Now())
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs("550e8400-e29b-41d4-a716-446655440000", models.NormalPurchase, decimal.NewFromFloat(-123.45)).
					WillReturnRows(txRows)
			},
			wantErr:      false,
			expectedSign: -1,
		},
		{
			name:          "successful credit voucher transaction",
			accountID:     "550e8400-e29b-41d4-a716-446655440000",
			operationType: models.CreditVoucher,
			amount:        decimal.NewFromFloat(100.00),
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "document_number", "created_at", "updated_at"}).
					AddRow("550e8400-e29b-41d4-a716-446655440000", "12345678900", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnRows(rows)

				txRows := sqlmock.NewRows([]string{"id", "account_id", "amount", "event_date", "operation_type", "created_at"}).
					AddRow("tx-id", "550e8400-e29b-41d4-a716-446655440000", decimal.NewFromFloat(100.00), time.Now(), "CREDIT_VOUCHER", time.Now())
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs("550e8400-e29b-41d4-a716-446655440000", models.CreditVoucher, decimal.NewFromFloat(100.00)).
					WillReturnRows(txRows)
			},
			wantErr:      false,
			expectedSign: 1,
		},
		{
			name:          "invalid operation type",
			accountID:     "550e8400-e29b-41d4-a716-446655440000",
			operationType: "INVALID_TYPE",
			amount:        decimal.NewFromFloat(123.45),
			mockSetup:     func() {},
			wantErr:       true,
			errContains:   "invalid operation type",
		},
		{
			name:          "invalid amount - zero",
			accountID:     "550e8400-e29b-41d4-a716-446655440000",
			operationType: models.NormalPurchase,
			amount:        decimal.Zero,
			mockSetup:     func() {},
			wantErr:       true,
			errContains:   "must be positive",
		},
		{
			name:          "invalid amount - negative",
			accountID:     "550e8400-e29b-41d4-a716-446655440000",
			operationType: models.NormalPurchase,
			amount:        decimal.NewFromFloat(-50.00),
			mockSetup:     func() {},
			wantErr:       true,
			errContains:   "must be positive",
		},
		{
			name:          "account not found",
			accountID:     "non-existent",
			operationType: models.NormalPurchase,
			amount:        decimal.NewFromFloat(123.45),
			mockSetup: func() {
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id").
					WithArgs("non-existent").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr:     true,
			errContains: "account not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			transaction, err := service.CreateTransaction(ctx, tt.accountID, tt.operationType, tt.amount)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, transaction)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, transaction)
				assert.Equal(t, tt.accountID, transaction.AccountID)
				assert.Equal(t, tt.operationType, transaction.OperationType)
				// Check amount sign
				if tt.expectedSign < 0 {
					assert.True(t, transaction.Amount.IsNegative(), "Amount should be negative")
				} else {
					assert.True(t, transaction.Amount.IsPositive(), "Amount should be positive")
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTransactionService_AdjustAmountSign(t *testing.T) {
	service := &TransactionService{}

	tests := []struct {
		name          string
		operationType models.OperationType
		amount        decimal.Decimal
		expected      decimal.Decimal
	}{
		{
			name:          "normal purchase - positive to negative",
			operationType: models.NormalPurchase,
			amount:        decimal.NewFromFloat(100.00),
			expected:      decimal.NewFromFloat(-100.00),
		},
		{
			name:          "withdrawal - positive to negative",
			operationType: models.Withdrawal,
			amount:        decimal.NewFromFloat(50.00),
			expected:      decimal.NewFromFloat(-50.00),
		},
		{
			name:          "credit voucher - stays positive",
			operationType: models.CreditVoucher,
			amount:        decimal.NewFromFloat(200.00),
			expected:      decimal.NewFromFloat(200.00),
		},
		{
			name:          "purchase with installments - positive to negative",
			operationType: models.PurchaseWithInstallments,
			amount:        decimal.NewFromFloat(300.00),
			expected:      decimal.NewFromFloat(-300.00),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.adjustAmountSign(tt.operationType, tt.amount)
			assert.True(t, tt.expected.Equal(result), "expected %s but got %s", tt.expected.String(), result.String())
		})
	}
}
