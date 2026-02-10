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

var noteValidate = validator.New()

// GetNotes returns all notes
func GetNotes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var notes []models.Note
		cursor, err := getNoteCollection().Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing notes"})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &notes); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, notes)
	}
}

// GetNoteByID returns a single note by ID
func GetNoteByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		noteId := c.Param("note_id")
		var note models.Note

		err := getNoteCollection().FindOne(ctx, bson.M{"note_id": noteId}).Decode(&note)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the note"})
			return
		}

		c.JSON(http.StatusOK, note)
	}
}

// CreateNote creates a new note
func CreateNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var note models.Note
		if err := c.BindJSON(&note); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := noteValidate.Struct(note)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		note.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		note.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		note.ID = primitive.NewObjectID()
		note.NoteID = note.ID.Hex()

		result, insertErr := getNoteCollection().InsertOne(ctx, note)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "note was not created"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// UpdateNote updates an existing note
func UpdateNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var note models.Note
		noteId := c.Param("note_id")

		if err := c.BindJSON(&note); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if note.Text != "" {
			updateObj = append(updateObj, bson.E{Key: "text", Value: note.Text})
		}

		if note.Title != "" {
			updateObj = append(updateObj, bson.E{Key: "title", Value: note.Title})
		}

		note.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: note.UpdatedAt})

		upsert := true
		filter := bson.M{"note_id": noteId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := getNoteCollection().UpdateOne(
			ctx,
			filter,
			bson.D{{Key: "$set", Value: updateObj}},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "note update failed"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetAllNotes returns all notes
func GetAllNotes() gin.HandlerFunc {
	return GetNotes()
}
