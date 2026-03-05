package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type OperationType string

const (
	NormalPurchase           OperationType = "NORMAL_PURCHASE"
	PurchaseWithInstallments OperationType = "PURCHASE_WITH_INSTALLMENTS"
	Withdrawal               OperationType = "WITHDRAWAL"
	CreditVoucher            OperationType = "CREDIT_VOUCHER"
)

type Transaction struct {
	ID            string          `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	AccountID     string          `json:"account_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Amount        decimal.Decimal `json:"amount" example:"-123.45"`
	EventDate     time.Time       `json:"event_date" example:"2026-03-05T12:00:00Z"`
	OperationType OperationType   `json:"operation_type" example:"NORMAL_PURCHASE"`
	CreatedAt     time.Time       `json:"created_at" example:"2026-03-05T12:00:00Z"`
}

type CreateTransactionRequest struct {
	AccountID     string          `json:"account_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440001"`
	OperationType OperationType   `json:"operation_type" binding:"required" example:"NORMAL_PURCHASE"`
	Amount        decimal.Decimal `json:"amount" binding:"required" example:"123.45"`
}

func (o OperationType) IsValid() bool {
	switch o {
	case NormalPurchase, PurchaseWithInstallments, Withdrawal, CreditVoucher:
		return true
	}
	return false
}
