package routes

func NoteRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/notes", controllers.GetNotes())
	incomingRoutes.GET("/notes/:note_id", controllers.GetNote())
	incomingRoutes.POST("/notes", controllers.CreateNotes())
	incomingRoutes.PATCH("/notes/:note_id", controllers.UpdateNote())
}