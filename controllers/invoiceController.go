package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ali-adel-nour/restaurant-management/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var invoiceValidate = validator.New()

// GetInvoices returns all invoices
func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var invoices []models.Invoice
		cursor, err := getInvoiceCollection().Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing invoices"})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &invoices); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, invoices)
	}
}

// GetInvoiceByID returns a single invoice by ID
func GetInvoiceByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		invoiceId := c.Param("invoice_id")
		var invoice models.Invoice

		err := getInvoiceCollection().FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the invoice"})
			return
		}

		c.JSON(http.StatusOK, invoice)
	}
}

// CreateInvoice creates a new invoice
func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var invoice models.Invoice
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := invoiceValidate.Struct(invoice)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Verify order exists
		var order models.Order
		err := getOrderCollection().FindOne(ctx, bson.M{"order_id": invoice.OrderID}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order was not found"})
			return
		}

		invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()
		invoice.InvoiceID = invoice.ID.Hex()

		// Set payment due date (e.g., 30 days from now)
		if invoice.PaymentDue == nil {
			paymentDue := time.Now().Add(30 * 24 * time.Hour)
			invoice.PaymentDue = &paymentDue
		}

		result, insertErr := getInvoiceCollection().InsertOne(ctx, invoice)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invoice was not created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// UpdateInvoice updates an existing invoice
func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var invoice models.Invoice
		invoiceId := c.Param("invoice_id")

		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if invoice.PaymentMethod != nil {
			updateObj = append(updateObj, bson.E{Key: "payment_method", Value: invoice.PaymentMethod})
		}

		if invoice.PaymentStatus != nil {
			updateObj = append(updateObj, bson.E{Key: "payment_status", Value: invoice.PaymentStatus})
		}

		if invoice.PaymentDue != nil {
			updateObj = append(updateObj, bson.E{Key: "payment_due", Value: invoice.PaymentDue})
		}

		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: invoice.UpdatedAt})

		upsert := true
		filter := bson.M{"invoice_id": invoiceId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := getInvoiceCollection().UpdateOne(
			ctx,
			filter,
			bson.D{{Key: "$set", Value: updateObj}},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invoice update failed"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetAllInvoices returns all invoices
func GetAllInvoices() gin.HandlerFunc {
	return GetInvoices()
}
