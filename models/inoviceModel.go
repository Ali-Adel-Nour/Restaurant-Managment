package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Invoice represents an invoice for an order
type Invoice struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	InvoiceID     string             `json:"invoice_id"`
	OrderID       string             `json:"order_id"`
	PaymentMethod *string            `json:"payment_method" validate:"eq=CARD|eq=CASH|eq="`
	PaymentStatus *string            `json:"payment_status" validate:"required,eq=PENDING|eq=PAID"`
	PaymentDue    *time.Time         `json:"payment_due"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}
