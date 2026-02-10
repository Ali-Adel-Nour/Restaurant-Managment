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

var tableValidate = validator.New()

// GetTables returns all tables
func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var tables []models.Table
		cursor, err := getTableCollection().Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing tables"})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &tables); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, tables)
	}
}

// GetTableByID returns a single table by ID
func GetTableByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tableId := c.Param("table_id")
		var table models.Table

		err := getTableCollection().FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the table"})
			return
		}

		c.JSON(http.StatusOK, table)
	}
}

// CreateTable creates a new table
func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table
		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := tableValidate.Struct(table)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		table.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.ID = primitive.NewObjectID()
		table.TableID = table.ID.Hex()

		result, insertErr := getTableCollection().InsertOne(ctx, table)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "table was not created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// UpdateTable updates an existing table
func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table
		tableId := c.Param("table_id")

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if table.NumberOfGuests != nil {
			updateObj = append(updateObj, bson.E{Key: "number_of_guests", Value: table.NumberOfGuests})
		}

		if table.TableNumber != nil {
			updateObj = append(updateObj, bson.E{Key: "table_number", Value: table.TableNumber})
		}

		table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: table.UpdatedAt})

		upsert := true
		filter := bson.M{"table_id": tableId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := getTableCollection().UpdateOne(
			ctx,
			filter,
			bson.D{{Key: "$set", Value: updateObj}},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "table update failed"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetAllTables returns all tables
func GetAllTables() gin.HandlerFunc {
	return GetTables()
}
