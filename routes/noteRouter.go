package routes

import (
	controller "github.com/ali-adel-nour/restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func NoteRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/notes", controller.GetAllNotes())
	incomingRoutes.GET("/notes/:note_id", controller.GetNoteByID())
	incomingRoutes.POST("/notes", controller.CreateNote())
	incomingRoutes.PATCH("/notes/:note_id", controller.UpdateNote())
}
