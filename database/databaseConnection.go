package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client is the MongoDB client instance
var Client *mongo.Client

// DB is the database instance
var DB *mongo.Database

// ConnectDB initializes the MongoDB connection with best practices
func ConnectDB() {
	// Get connection string from environment
	MongoDB := os.Getenv("MONGODB_URL")
	if MongoDB == "" {
		MongoDB = "mongodb://localhost:27017"
	}

	// Get database name from environment
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "restaurant"
	}

	// Set client options with best practices
	clientOptions := options.Client().
		ApplyURI(MongoDB).
		SetMaxPoolSize(50).                        // Maximum connections in pool
		SetMinPoolSize(10).                        // Minimum connections in pool
		SetMaxConnIdleTime(30 * time.Second).      // Close idle connections after 30s
		SetConnectTimeout(10 * time.Second).       // Connection timeout
		SetServerSelectionTimeout(5 * time.Second) // Server selection timeout

	// Create context with timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping to verify connection
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	fmt.Println("✅ Connected to MongoDB successfully!")
	Client = client
	DB = client.Database(dbName)
}

// GetCollection returns a collection from the database
// When called with a nil client, it will use the global Client once initialized
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "restaurant"
	}

	// Use global Client if available
	if Client != nil {
		return Client.Database(dbName).Collection(collectionName)
	}

	// Use provided client if available
	if client != nil {
		return client.Database(dbName).Collection(collectionName)
	}

	// Return nil - will be properly initialized after ConnectDB
	return nil
}

// OpenCollection gets a collection from the connected database
// Use this to get collections after ConnectDB() has been called
func OpenCollection(collectionName string) *mongo.Collection {
	if DB == nil {
		// If DB is not initialized yet, try to use Client directly
		if Client == nil {
			log.Fatal("Database not initialized. Call ConnectDB() first.")
		}
		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "restaurant"
		}
		return Client.Database(dbName).Collection(collectionName)
	}
	return DB.Collection(collectionName)
}

// DisconnectDB gracefully closes the MongoDB connection
func DisconnectDB() error {
	if Client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := Client.Disconnect(ctx)
	if err != nil {
		return err
	}

	fmt.Println("✅ Disconnected from MongoDB")
	return nil
}

// DBinstance returns the MongoDB client (kept for backward compatibility)
func DBinstance() *mongo.Client {
	return Client
}
