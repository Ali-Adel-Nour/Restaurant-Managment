package controller

import (
	"github.com/ali-adel-nour/restaurant-management/database"
	"go.mongodb.org/mongo-driver/mongo"
)

// Collection getter functions
func getUserCollection() *mongo.Collection {
	return database.Collections.Users
}

func getFoodCollection() *mongo.Collection {
	return database.Collections.Foods
}

func getMenuCollection() *mongo.Collection {
	return database.Collections.Menus
}

func getTableCollection() *mongo.Collection {
	return database.Collections.Tables
}

func getOrderCollection() *mongo.Collection {
	return database.Collections.Orders
}

func getOrderItemCollection() *mongo.Collection {
	return database.Collections.OrderItems
}

func getInvoiceCollection() *mongo.Collection {
	return database.Collections.Invoices
}

func getNoteCollection() *mongo.Collection {
	return database.Collections.Notes
}
