package routes

import (
	controller "github.com/ali-adel-nour/restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/foods", controller.GetAllFoods())
	incomingRoutes.GET("/foods/:food_id", controller.GetFoodByID())
	incomingRoutes.POST("/foods", controller.CreateFood())
	incomingRoutes.PATCH("/foods/:food_id", controller.UpdateFood())
}
