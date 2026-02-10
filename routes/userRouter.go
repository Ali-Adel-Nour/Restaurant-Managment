package routes

import (
	controller "github.com/ali-adel-nour/restaurant-management/controllers"
	"github.com/ali-adel-nour/restaurant-management/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", middleware.Authentication(), controller.GetAllUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.POST("/users/signup", controller.SignUp())
	incomingRoutes.POST("/users/login", controller.Login())
	incomingRoutes.POST("/users/logout", middleware.Authentication(), controller.Logout())
}
