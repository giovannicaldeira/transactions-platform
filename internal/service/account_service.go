package service

import (
	"context"
	"fmt"

	"github.com/transactions-platform/internal/models"
	"github.com/transactions-platform/internal/repository"
)

type AccountService struct {
	accountRepo *repository.AccountRepository
}

func NewAccountService(accountRepo *repository.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

// CreateAccount creates a new account with business logic validation
func (s *AccountService) CreateAccount(ctx context.Context, documentNumber string) (*models.Account, error) {
	// Check if account with this document number already exists
	existing, err := s.accountRepo.GetByDocumentNumber(ctx, documentNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing account: %w", err)
	}

	if existing != nil {
		return nil, fmt.Errorf("account with document number %s already exists", documentNumber)
	}

	// Create the account
	account, err := s.accountRepo.Create(ctx, documentNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

// GetAccountByID retrieves an account by ID
func (s *AccountService) GetAccountByID(ctx context.Context, id string) (*models.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "account not found" {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

// GetAccountByDocumentNumber retrieves an account by document number
func (s *AccountService) GetAccountByDocumentNumber(ctx context.Context, documentNumber string) (*models.Account, error) {
	return s.accountRepo.GetByDocumentNumber(ctx, documentNumber)
}
