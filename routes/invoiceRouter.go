package routes

import (
	"github.com/ahmedabzk/restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/invoices", controllers.GetInvoices())
	incomingRoutes.GET("/invoices/:invoice_id", controllers.GetInvoice())
	incomingRoutes.POST("/invoices", controllers.CreateInvoices())
	incomingRoutes.PATCH("/invoices/:invoice_id", controllers.UpdateInvoice())
}