package routes

import (
	controller "github.com/ali-adel-nour/restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/menus", controller.GetAllMenus())
	incomingRoutes.GET("/menus/:menu_id", controller.GetMenuByID())
	incomingRoutes.POST("/menus", controller.CreateMenu())
	incomingRoutes.PATCH("/menus/:menu_id", controller.UpdateMenu())
}
