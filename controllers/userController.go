package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ali-adel-nour/restaurant-management/helpers"
	"github.com/ali-adel-nour/restaurant-management/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

// HashPassword hashes a password using bcrypt
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

// VerifyPassword verifies if the provided password matches the hashed password
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "login or password is incorrect"
		check = false
	}

	return check, msg
}

// GetUsers returns all users (admin only)
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var users []models.User
		cursor, err := getUserCollection().Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing users"})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &users); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, users)
	}
}

// GetUser returns a single user by ID
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		err := getUserCollection().FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching user"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// SignUp creates a new user account
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate user input
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if email already exists
		count, err := getUserCollection().CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for email"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}

		// Hash password
		password := HashPassword(*user.Password)
		user.Password = &password

		// Check if phone already exists
		count, err = getUserCollection().CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for phone"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "phone number already exists"})
			return
		}

		// Set timestamps and IDs
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()

		// Generate tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, user.UserID)
		user.Token = &token
		user.RefreshToken = &refreshToken

		// Insert user
		_, insertErr := getUserCollection().InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user was not created"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// Login authenticates a user
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find user by email
		err := getUserCollection().FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "login or password is incorrect"})
			return
		}

		// Verify password
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		// Generate new tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserID)

		// Update tokens in database
		helpers.UpdateAllTokens(token, refreshToken, foundUser.UserID)

		// Return user data with new tokens
		foundUser.Token = &token
		foundUser.RefreshToken = &refreshToken

		c.JSON(http.StatusOK, foundUser)
	}
}

// CreateUser creates a new user (admin function)
func CreateUser() gin.HandlerFunc {
	return SignUp()
}

// GetAllUsers returns all users
func GetAllUsers() gin.HandlerFunc {
	return GetUsers()
}

// Logout logs out a user
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	}
}
