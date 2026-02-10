package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Quantity    *int               `json:"quantity" validate:"required,eq=1|eq=2|eq=3|eq=4|eq=5"`
	UnitPrice   *float64           `json:"unit_price" validate:"required"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	FoodID      *string            `json:"food_id" validate:"required"`
	OrderItemID string             `json:"order_item_id"`
	OrderID     string             `json:"order_id" validate:"required"`
}
