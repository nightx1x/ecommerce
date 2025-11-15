package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id,omitempty"`
	ProductID uuid.UUID `db:"product_id" json:"product_id,omitempty"`
	Quantity  int       `db:"quantity" json:"quantity"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CartWithProduct struct {
	Cart
	ProductName  string  `db:"product_name" json:"product_name"`
	ProductPrice float64 `db:"product_price" json:"product_price"`
	ProductStock int     `db:"product_stock" json:"product_stock"`
}
