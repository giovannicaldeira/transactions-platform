package models

import "time"

type Account struct {
	ID             string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	DocumentNumber string    `json:"document_number" example:"12345678900"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
