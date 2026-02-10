package routes

import (
	controller "github.com/ali-adel-nour/restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/invoices", controller.GetAllInvoices())
	incomingRoutes.GET("/invoices/:invoice_id", controller.GetInvoiceByID())
	incomingRoutes.POST("/invoices", controller.CreateInvoice())
	incomingRoutes.PATCH("/invoices/:invoice_id", controller.UpdateInvoice())

}
