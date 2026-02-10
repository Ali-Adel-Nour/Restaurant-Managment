package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ali-adel-nour/restaurant-management/database"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SignedDetails represents the JWT token claims
type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UID       string
	jwt.StandardClaims
}

// GetSecretKey returns the JWT secret key from environment
func GetSecretKey() string {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		secret = "your-secret-key-change-this-in-production"
	}
	return secret
}

// GenerateAllTokens generates both access and refresh tokens
func GenerateAllTokens(email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {
	// Access token claims (valid for 24 hours)
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UID:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24).Unix(),
		},
	}

	// Refresh token claims (valid for 7 days)
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 168).Unix(),
		},
	}

	// Generate tokens
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(GetSecretKey()))
	if err != nil {
		log.Panic(err)
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(GetSecretKey()))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

// UpdateAllTokens updates the user's tokens in the database
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: updatedAt})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := database.Collections.Users.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)

	if err != nil {
		log.Panic(err)
		return
	}
}

// ValidateToken validates the JWT token and returns the claims
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(GetSecretKey()), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}

	return claims, msg
}
