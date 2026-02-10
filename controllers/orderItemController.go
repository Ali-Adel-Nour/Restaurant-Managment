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

var orderItemValidate = validator.New()

// GetOrderItems returns all order items
func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItems []models.OrderItem
		cursor, err := getOrderItemCollection().Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing order items"})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &orderItems); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, orderItems)
	}
}

// GetOrderItemByID returns a single order item by ID
func GetOrderItemByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderItemId := c.Param("orderItem_id")
		var orderItem models.OrderItem

		err := getOrderItemCollection().FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the order item"})
			return
		}

		c.JSON(http.StatusOK, orderItem)
	}
}

// GetOrderItemsByOrderID returns all order items for a specific order
func GetOrderItemsByOrderID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderId := c.Param("order_id")

		var orderItems []models.OrderItem
		cursor, err := getOrderItemCollection().Find(ctx, bson.M{"order_id": orderId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching order items"})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &orderItems); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, orderItems)
	}
}

// CreateOrderItem creates a new order item
func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItem models.OrderItem
		if err := c.BindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := orderItemValidate.Struct(orderItem)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Verify order exists
		var order models.Order
		err := getOrderCollection().FindOne(ctx, bson.M{"order_id": orderItem.OrderID}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order was not found"})
			return
		}

		// Verify food exists
		var food models.Food
		err = getFoodCollection().FindOne(ctx, bson.M{"food_id": orderItem.FoodID}).Decode(&food)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "food item was not found"})
			return
		}

		orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.ID = primitive.NewObjectID()
		orderItem.OrderItemID = orderItem.ID.Hex()

		// Set unit price from food price if not provided
		if orderItem.UnitPrice == nil {
			orderItem.UnitPrice = food.Price
		}

		result, insertErr := getOrderItemCollection().InsertOne(ctx, orderItem)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order item was not created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// UpdateOrderItem updates an existing order item
func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItem models.OrderItem
		orderItemId := c.Param("orderItem_id")

		if err := c.BindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{Key: "quantity", Value: orderItem.Quantity})
		}

		if orderItem.UnitPrice != nil {
			updateObj = append(updateObj, bson.E{Key: "unit_price", Value: orderItem.UnitPrice})
		}

		if orderItem.FoodID != nil {
			var food models.Food
			err := getFoodCollection().FindOne(ctx, bson.M{"food_id": orderItem.FoodID}).Decode(&food)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "food item was not found"})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "food_id", Value: orderItem.FoodID})
		}

		orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: orderItem.UpdatedAt})

		upsert := true
		filter := bson.M{"order_item_id": orderItemId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := getOrderItemCollection().UpdateOne(
			ctx,
			filter,
			bson.D{{Key: "$set", Value: updateObj}},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order item update failed"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
