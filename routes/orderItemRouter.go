package routes

import "github.com/ahmedabzk/restaurant-management/controllers"

func OrderItemRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/orderItems", controllers.GetOrderItems())
	incomingRoutes.GET("/orderItems/:orderItem_id", controllers.GetOrderItem())
	incomingRoutes.GET("/orderItems-order/:order_id", controllers.GetOrderItemsByOrder())
	incomingRoutes.POST("/orderItems", controllers.CreateOrderItems())
	incomingRoutes.PATCH("/orderItems/:orderItem_id", controllers.UpdateTable())
}