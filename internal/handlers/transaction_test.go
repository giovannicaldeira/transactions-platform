package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/transactions-platform/internal/dto"
	"github.com/transactions-platform/internal/models"
	"github.com/transactions-platform/internal/repository"
	"github.com/transactions-platform/internal/service"
)

func setupTransactionTestRouter() (*gin.Engine, sqlmock.Sqlmock, *sql.DB) {
	gin.SetMode(gin.TestMode)

	db, mock, _ := sqlmock.New()
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo)
	handler := NewTransactionHandler(transactionService)

	router := gin.New()
	router.POST("/transactions", handler.CreateTransaction)

	return router, mock, db
}

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    any
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful purchase transaction",
			requestBody: dto.CreateTransactionRequest{
				AccountID:     "550e8400-e29b-41d4-a716-446655440000",
				OperationType: models.NormalPurchase,
				Amount:        decimal.NewFromFloat(123.45),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Account exists check
				accountRows := sqlmock.NewRows([]string{"id", "document_number", "created_at", "updated_at"}).
					AddRow("550e8400-e29b-41d4-a716-446655440000", "12345678900", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnRows(accountRows)

				// Create transaction (amount should be negative for purchase)
				txRows := sqlmock.NewRows([]string{"id", "account_id", "amount", "event_date", "operation_type", "created_at"}).
					AddRow("tx-id", "550e8400-e29b-41d4-a716-446655440000", decimal.NewFromFloat(-123.45), time.Now(), "NORMAL_PURCHASE", time.Now())
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs("550e8400-e29b-41d4-a716-446655440000", models.NormalPurchase, decimal.NewFromFloat(-123.45)).
					WillReturnRows(txRows)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var transaction models.Transaction
				err := json.Unmarshal(w.Body.Bytes(), &transaction)
				require.NoError(t, err)
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", transaction.AccountID)
				assert.Equal(t, models.NormalPurchase, transaction.OperationType)
				assert.True(t, transaction.Amount.IsNegative(), "Purchase amount should be negative")
			},
		},
		{
			name: "successful credit voucher transaction",
			requestBody: dto.CreateTransactionRequest{
				AccountID:     "550e8400-e29b-41d4-a716-446655440000",
				OperationType: models.CreditVoucher,
				Amount:        decimal.NewFromFloat(100.00),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				accountRows := sqlmock.NewRows([]string{"id", "document_number", "created_at", "updated_at"}).
					AddRow("550e8400-e29b-41d4-a716-446655440000", "12345678900", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnRows(accountRows)

				txRows := sqlmock.NewRows([]string{"id", "account_id", "amount", "event_date", "operation_type", "created_at"}).
					AddRow("tx-id", "550e8400-e29b-41d4-a716-446655440000", decimal.NewFromFloat(100.00), time.Now(), "CREDIT_VOUCHER", time.Now())
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs("550e8400-e29b-41d4-a716-446655440000", models.CreditVoucher, decimal.NewFromFloat(100.00)).
					WillReturnRows(txRows)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var transaction models.Transaction
				err := json.Unmarshal(w.Body.Bytes(), &transaction)
				require.NoError(t, err)
				assert.Equal(t, models.CreditVoucher, transaction.OperationType)
				assert.True(t, transaction.Amount.IsPositive(), "Credit amount should be positive")
			},
		},
		{
			name:           "invalid request body",
			requestBody:    map[string]string{"invalid": "data"},
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "error")
			},
		},
		{
			name: "invalid operation type",
			requestBody: dto.CreateTransactionRequest{
				AccountID:     "550e8400-e29b-41d4-a716-446655440000",
				OperationType: "INVALID_TYPE",
				Amount:        decimal.NewFromFloat(123.45),
			},
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "invalid operation type")
			},
		},
		{
			name: "invalid amount - zero",
			requestBody: dto.CreateTransactionRequest{
				AccountID:     "550e8400-e29b-41d4-a716-446655440000",
				OperationType: models.NormalPurchase,
				Amount:        decimal.Zero,
			},
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "must be positive")
			},
		},
		{
			name: "invalid amount - negative",
			requestBody: dto.CreateTransactionRequest{
				AccountID:     "550e8400-e29b-41d4-a716-446655440000",
				OperationType: models.NormalPurchase,
				Amount:        decimal.NewFromFloat(-50.00),
			},
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "must be positive")
			},
		},
		{
			name: "account not found",
			requestBody: dto.CreateTransactionRequest{
				AccountID:     "non-existent",
				OperationType: models.NormalPurchase,
				Amount:        decimal.NewFromFloat(123.45),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id").
					WithArgs("non-existent").
					WillReturnError(sql.ErrNoRows)
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "not found")
			},
		},
		{
			name: "successful withdrawal transaction",
			requestBody: dto.CreateTransactionRequest{
				AccountID:     "550e8400-e29b-41d4-a716-446655440000",
				OperationType: models.Withdrawal,
				Amount:        decimal.NewFromFloat(50.00),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				accountRows := sqlmock.NewRows([]string{"id", "document_number", "created_at", "updated_at"}).
					AddRow("550e8400-e29b-41d4-a716-446655440000", "12345678900", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnRows(accountRows)

				txRows := sqlmock.NewRows([]string{"id", "account_id", "amount", "event_date", "operation_type", "created_at"}).
					AddRow("tx-id", "550e8400-e29b-41d4-a716-446655440000", decimal.NewFromFloat(-50.00), time.Now(), "WITHDRAWAL", time.Now())
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs("550e8400-e29b-41d4-a716-446655440000", models.Withdrawal, decimal.NewFromFloat(-50.00)).
					WillReturnRows(txRows)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var transaction models.Transaction
				err := json.Unmarshal(w.Body.Bytes(), &transaction)
				require.NoError(t, err)
				assert.Equal(t, models.Withdrawal, transaction.OperationType)
				assert.True(t, transaction.Amount.IsNegative(), "Withdrawal amount should be negative")
			},
		},
		{
			name: "successful installment purchase transaction",
			requestBody: dto.CreateTransactionRequest{
				AccountID:     "550e8400-e29b-41d4-a716-446655440000",
				OperationType: models.PurchaseWithInstallments,
				Amount:        decimal.NewFromFloat(300.00),
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				accountRows := sqlmock.NewRows([]string{"id", "document_number", "created_at", "updated_at"}).
					AddRow("550e8400-e29b-41d4-a716-446655440000", "12345678900", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnRows(accountRows)

				txRows := sqlmock.NewRows([]string{"id", "account_id", "amount", "event_date", "operation_type", "created_at"}).
					AddRow("tx-id", "550e8400-e29b-41d4-a716-446655440000", decimal.NewFromFloat(-300.00), time.Now(), "PURCHASE_WITH_INSTALLMENTS", time.Now())
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs("550e8400-e29b-41d4-a716-446655440000", models.PurchaseWithInstallments, decimal.NewFromFloat(-300.00)).
					WillReturnRows(txRows)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var transaction models.Transaction
				err := json.Unmarshal(w.Body.Bytes(), &transaction)
				require.NoError(t, err)
				assert.Equal(t, models.PurchaseWithInstallments, transaction.OperationType)
				assert.True(t, transaction.Amount.IsNegative(), "Installment purchase amount should be negative")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mock, db := setupTransactionTestRouter()
			defer db.Close()

			tt.mockSetup(mock)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
