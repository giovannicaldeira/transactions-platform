package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/transactions-platform/internal/logger"
	"github.com/transactions-platform/internal/models"
	"github.com/transactions-platform/internal/service"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// CreateTransaction godoc
// @Summary      Create a new transaction
// @Description  Creates a new transaction for an account
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        request body models.CreateTransactionRequest true "Transaction creation request"
// @Success      201  {object}  models.Transaction
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      422  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req models.CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request body for transaction creation").
			Err(err).
			Str("ip", c.ClientIP()).
			Send()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	logger.Info("Creating new transaction").
		Str("account_id", req.AccountID).
		Str("operation_type", string(req.OperationType)).
		Str("amount", req.Amount.String()).
		Send()

	// Call service to create transaction
	transaction, err := h.service.CreateTransaction(c.Request.Context(), req.AccountID, req.OperationType, req.Amount)
	if err != nil {
		// Determine appropriate status code based on error
		if strings.Contains(err.Error(), "not found") {
			logger.Warn("Transaction creation failed - account not found").
				Str("account_id", req.AccountID).
				Send()
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "invalid operation type") {
			logger.Warn("Transaction creation failed - invalid operation type").
				Str("operation_type", string(req.OperationType)).
				Send()
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "must be positive") {
			logger.Warn("Transaction creation failed - invalid amount").
				Str("amount", req.Amount.String()).
				Send()
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		logger.Error("Failed to create transaction").
			Err(err).
			Str("account_id", req.AccountID).
			Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	logger.Info("Transaction created successfully").
		Str("transaction_id", transaction.ID).
		Str("account_id", transaction.AccountID).
		Str("operation_type", string(transaction.OperationType)).
		Str("amount", transaction.Amount.String()).
		Send()

	c.JSON(http.StatusCreated, transaction)
}

