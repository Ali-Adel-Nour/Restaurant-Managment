package main

import (
	"os"

	"github.com/ali-adel-nour/restaurant-management/database"
	"github.com/ali-adel-nour/restaurant-management/middleware"
	"github.com/ali-adel-nour/restaurant-management/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to MongoDB
	database.ConnectDB()

	// Initialize collections
	database.InitCollections()

	// Create Gin router
	router := gin.New()
	router.Use(gin.Logger())

	// Public routes (no authentication required)
	routes.UserRoutes(router)

	// Protected routes (authentication required)
	router.Use(middleware.Authentication())
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.OrderRoutes(router)
	routes.TableRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)
	routes.NoteRoutes(router)

	// Start server
	router.Run(":" + port)
}
