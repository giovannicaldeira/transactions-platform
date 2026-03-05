package models

import "time"

type Account struct {
	ID             string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	DocumentNumber string    `json:"document_number" example:"12345678900"`
	CreatedAt      time.Time `json:"created_at" example:"2026-03-04T12:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2026-03-04T12:00:00Z"`
}

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" binding:"required" example:"12345678900"`
}
