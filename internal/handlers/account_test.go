package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/transactions-platform/internal/models"
	"github.com/transactions-platform/internal/repository"
)

func setupTestRouter() (*gin.Engine, sqlmock.Sqlmock, *sql.DB) {
	gin.SetMode(gin.TestMode)

	db, mock, _ := sqlmock.New()
	repo := repository.NewAccountRepository(db)
	handler := NewAccountHandler(repo)

	router := gin.New()
	router.POST("/accounts", handler.CreateAccount)
	router.GET("/accounts/:id", handler.GetAccount)

	return router, mock, db
}

func TestAccountHandler_CreateAccount(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful account creation",
			requestBody: models.CreateAccountRequest{
				DocumentNumber: "12345678900",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Check for existing account
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
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var account models.Account
				err := json.Unmarshal(w.Body.Bytes(), &account)
				require.NoError(t, err)
				assert.Equal(t, "12345678900", account.DocumentNumber)
				assert.NotEmpty(t, account.ID)
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
			name: "duplicate document number",
			requestBody: models.CreateAccountRequest{
				DocumentNumber: "12345678900",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Return existing account
				rows := sqlmock.NewRows([]string{"id", "document_number", "created_at", "updated_at"}).
					AddRow("550e8400-e29b-41d4-a716-446655440000", "12345678900", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE document_number").
					WithArgs("12345678900").
					WillReturnRows(rows)
			},
			expectedStatus: http.StatusConflict,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "already exists")
			},
		},
		{
			name: "database error on check",
			requestBody: models.CreateAccountRequest{
				DocumentNumber: "12345678900",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE document_number").
					WithArgs("12345678900").
					WillReturnError(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "error")
			},
		},
		{
			name: "database error on create",
			requestBody: models.CreateAccountRequest{
				DocumentNumber: "12345678900",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE document_number").
					WithArgs("12345678900").
					WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery("INSERT INTO accounts").
					WithArgs("12345678900").
					WillReturnError(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "Failed to create account")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mock, db := setupTestRouter()
			defer db.Close()

			tt.mockSetup(mock)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
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

func TestAccountHandler_GetAccount(t *testing.T) {
	tests := []struct {
		name           string
		accountID      string
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "account found",
			accountID: "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "document_number", "created_at", "updated_at"}).
					AddRow("550e8400-e29b-41d4-a716-446655440000", "12345678900", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnRows(rows)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var account models.Account
				err := json.Unmarshal(w.Body.Bytes(), &account)
				require.NoError(t, err)
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", account.ID)
				assert.Equal(t, "12345678900", account.DocumentNumber)
			},
		},
		{
			name:      "account not found",
			accountID: "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnError(sql.ErrNoRows)
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "Account not found")
			},
		},
		{
			name:      "database error",
			accountID: "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM accounts WHERE id").
					WithArgs("550e8400-e29b-41d4-a716-446655440000").
					WillReturnError(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "Failed to get account")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mock, db := setupTestRouter()
			defer db.Close()

			tt.mockSetup(mock)

			req := httptest.NewRequest(http.MethodGet, "/accounts/"+tt.accountID, nil)
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

func TestAccountHandler_GetAccount_InvalidID(t *testing.T) {
	router, _, db := setupTestRouter()
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/accounts/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return 404 because route doesn't match
	assert.Equal(t, http.StatusNotFound, w.Code)
}
