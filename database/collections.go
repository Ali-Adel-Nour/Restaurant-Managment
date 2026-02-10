package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Collections holds all database collections
var Collections struct {
	Users      *mongo.Collection
	Foods      *mongo.Collection
	Menus      *mongo.Collection
	Tables     *mongo.Collection
	Orders     *mongo.Collection
	OrderItems *mongo.Collection
	Invoices   *mongo.Collection
	Notes      *mongo.Collection
}

// InitCollections initializes all database collections
// Call this after ConnectDB()
func InitCollections() {
	Collections.Users = OpenCollection("users")
	Collections.Foods = OpenCollection("foods")
	Collections.Menus = OpenCollection("menus")
	Collections.Tables = OpenCollection("tables")
	Collections.Orders = OpenCollection("orders")
	Collections.OrderItems = OpenCollection("orderItems")
	Collections.Invoices = OpenCollection("invoices")
	Collections.Notes = OpenCollection("notes")
}
