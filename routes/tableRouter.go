package routes

import (
	"github.com/ahmedabzk/restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func TableRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/tables", controllers.GetTables())
	incomingRoutes.GET("/tables/:table_id", controllers.GetTable())
	incomingRoutes.POST("/tables", controllers.CreateTables())
	incomingRoutes.PATCH("/tables/:table_id", controllers.UpdateTable())
}
