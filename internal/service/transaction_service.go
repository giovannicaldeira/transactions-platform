package service

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/transactions-platform/internal/models"
	"github.com/transactions-platform/internal/repository"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	accountRepo     *repository.AccountRepository
}

func NewTransactionService(transactionRepo *repository.TransactionRepository, accountRepo *repository.AccountRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

// CreateTransaction creates a new transaction with business logic validation
func (s *TransactionService) CreateTransaction(ctx context.Context, accountID string, operationType models.OperationType, amount decimal.Decimal) (*models.Transaction, error) {
	// Validate operation type
	if !operationType.IsValid() {
		return nil, fmt.Errorf("invalid operation type: %s", operationType)
	}

	// Validate amount
	if err := s.validateAmount(amount); err != nil {
		return nil, err
	}

	// Check if account exists
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		if err.Error() == "account not found" {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to verify account: %w", err)
	}

	if account == nil {
		return nil, fmt.Errorf("account not found")
	}

	// Adjust amount sign based on operation type
	adjustedAmount := s.adjustAmountSign(operationType, amount)

	// Create the transaction
	transaction, err := s.transactionRepo.Create(ctx, accountID, operationType, adjustedAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}


// validateAmount validates that the amount is positive
func (s *TransactionService) validateAmount(amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("amount must be positive, got: %s", amount.String())
	}
	return nil
}

// adjustAmountSign adjusts the amount sign based on operation type
// Purchases and withdrawals should be negative (debits)
// Credit vouchers should be positive (credits)
func (s *TransactionService) adjustAmountSign(operationType models.OperationType, amount decimal.Decimal) decimal.Decimal {
	// Ensure amount is positive first (use absolute value)
	amount = amount.Abs()

	switch operationType {
	case models.NormalPurchase, models.PurchaseWithInstallments, models.Withdrawal:
		return amount.Neg() // Debit (negative)
	case models.CreditVoucher:
		return amount // Credit (positive)
	default:
		return amount
	}
}
