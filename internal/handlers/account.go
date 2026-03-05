package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/transactions-platform/internal/dto"
	"github.com/transactions-platform/internal/logger"
	_ "github.com/transactions-platform/internal/models" // imported for swagger docs
	"github.com/transactions-platform/internal/service"
)

type AccountHandler struct {
	service *service.AccountService
}

func NewAccountHandler(service *service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

// CreateAccount godoc
// @Summary      Create a new account
// @Description  Creates a new account with the provided document number
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateAccountRequest true "Account creation request"
// @Success      201  {object}  models.Account
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req dto.CreateAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request body for account creation").
			Err(err).
			Str("ip", c.ClientIP()).
			Send()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	logger.Info("Creating new account").
		Str("document_number", req.DocumentNumber).
		Send()

	// Call service to create account
	account, err := h.service.CreateAccount(c.Request.Context(), req.DocumentNumber)
	if err != nil {
		// Determine appropriate status code based on error
		if strings.Contains(err.Error(), "already exists") {
			logger.Warn("Account creation failed - duplicate document number").
				Str("document_number", req.DocumentNumber).
				Send()
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		logger.Error("Failed to create account").
			Err(err).
			Str("document_number", req.DocumentNumber).
			Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	logger.Info("Account created successfully").
		Str("account_id", account.ID).
		Str("document_number", account.DocumentNumber).
		Send()

	c.JSON(http.StatusCreated, account)
}

// GetAccount godoc
// @Summary      Get account by ID
// @Description  Retrieves an account by its ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id path string true "Account ID"
// @Success      200  {object}  models.Account
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /accounts/{id} [get]
func (h *AccountHandler) GetAccount(c *gin.Context) {
	id := c.Param("id")

	logger.Debug("Fetching account").Str("account_id", id).Send()

	account, err := h.service.GetAccountByID(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			logger.Debug("Account not found").Str("account_id", id).Send()
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		logger.Error("Failed to get account").
			Err(err).
			Str("account_id", id).
			Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account"})
		return
	}

	logger.Debug("Account fetched successfully").Str("account_id", id).Send()

	c.JSON(http.StatusOK, account)
}
