package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID         `json:"id"`
	UserID     uuid.UUID         `json:"user_id"`
	Items      map[uuid.UUID]int `json:"items"` // product ID to quantity
	TotalPrice float64           `json:"total_price"`
	Status     string            `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type OrderFilter struct {
	UserID *uuid.UUID
	Limit  int
	Offset int
}
