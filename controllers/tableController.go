package controllers

import (
	"github.com/ahmedabzk/restaurant-management/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func GetTable() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func CreateTables() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func UpdateTable() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		
	}
}