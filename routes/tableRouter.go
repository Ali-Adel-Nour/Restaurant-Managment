package routes

import (
	controller "github.com/ali-adel-nour/restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func TableRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/tables", controller.GetAllTables())
	incomingRoutes.GET("/tables/:table_id", controller.GetTableByID())
	incomingRoutes.POST("/tables", controller.CreateTable())
	incomingRoutes.PATCH("/tables/:table_id", controller.UpdateTable())
}
