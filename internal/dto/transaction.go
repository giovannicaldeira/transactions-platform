package dto

import (
	"github.com/shopspring/decimal"
	"github.com/transactions-platform/internal/models"
)

type CreateTransactionRequest struct {
	AccountID     string                `json:"account_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440001"`
	OperationType models.OperationType  `json:"operation_type" binding:"required" example:"NORMAL_PURCHASE"`
	Amount        decimal.Decimal       `json:"amount" binding:"required" example:"123.45"`
}
