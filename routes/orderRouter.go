package routes

import (
	"github.com/ahmedabzk/restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/orders", controllers.GetOrders())
	incomingRoutes.GET("/orders/:order_id", controllers.GetOrder())
	incomingRoutes.POST("/orders", controllers.CreateOrders())
	incomingRoutes.PATCH("/orders/:order_id", controllers.UpdateOrder())
}
