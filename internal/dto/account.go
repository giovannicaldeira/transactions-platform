package dto

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" binding:"required" example:"12345678900"`
}
