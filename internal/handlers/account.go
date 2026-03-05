package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/transactions-platform/internal/models"
	"github.com/transactions-platform/internal/repository"
)

type AccountHandler struct {
	repo *repository.AccountRepository
}

func NewAccountHandler(repo *repository.AccountRepository) *AccountHandler {
	return &AccountHandler{repo: repo}
}

// CreateAccount godoc
// @Summary      Create a new account
// @Description  Creates a new account with the provided document number
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        request body models.CreateAccountRequest true "Account creation request"
// @Success      201  {object}  models.Account
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req models.CreateAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Check if account with this document number already exists
	existing, err := h.repo.GetByDocumentNumber(c.Request.Context(), req.DocumentNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing account"})
		return
	}

	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Account with this document number already exists"})
		return
	}

	// Create the account
	account, err := h.repo.Create(c.Request.Context(), req.DocumentNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

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

	account, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "account not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account"})
		return
	}

	c.JSON(http.StatusOK, account)
}
