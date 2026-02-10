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

var orderValidate = validator.New()

// GetOrders returns all orders
func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orders []models.Order
		cursor, err := getOrderCollection().Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing orders"})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &orders); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, orders)
	}
}

// GetOrderByID returns a single order by ID
func GetOrderByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderId := c.Param("order_id")
		var order models.Order

		err := getOrderCollection().FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the order"})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

// CreateOrder creates a new order
func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var order models.Order
		var table models.Table

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := orderValidate.Struct(order)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Verify table exists
		if order.TableID != nil {
			err := getTableCollection().FindOne(ctx, bson.M{"table_id": order.TableID}).Decode(&table)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "table was not found"})
				return
			}
		}

		order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.ID = primitive.NewObjectID()
		order.OrderID = order.ID.Hex()

		result, insertErr := getOrderCollection().InsertOne(ctx, order)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order was not created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// UpdateOrder updates an existing order
func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var order models.Order
		orderId := c.Param("order_id")

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if order.TableID != nil {
			var table models.Table
			err := getTableCollection().FindOne(ctx, bson.M{"table_id": order.TableID}).Decode(&table)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "table was not found"})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "table_id", Value: order.TableID})
		}

		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: order.UpdatedAt})

		upsert := true
		filter := bson.M{"order_id": orderId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := getOrderCollection().UpdateOne(
			ctx,
			filter,
			bson.D{{Key: "$set", Value: updateObj}},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order update failed"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetAllOrders returns all orders
func GetAllOrders() gin.HandlerFunc {
	return GetOrders()
}
