package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/transactions-platform/internal/models"
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

	// Call service to create account
	account, err := h.service.CreateAccount(c.Request.Context(), req.DocumentNumber)
	if err != nil {
		// Determine appropriate status code based on error
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
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

	account, err := h.service.GetAccountByID(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account"})
		return
	}

	c.JSON(http.StatusOK, account)
}
