package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/ahmedabzk/restaurant-management/database"
	"github.com/ahmedabzk/restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func GetFoods() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}

func GetFood() gin.HandlerFunc{
	return func(c *gin.Context) {
		// use context.WithTimeout so that the API terminates after the specified time
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// get the id of the food item we looking for from the Params received
		foodId := c.Param("food_id")

		var food models.Food
		// look for the food item from the collection using the foodId
		err := foodCollection.FindOne(ctx, bson.M{"food_id":foodId}).Decode(&food)

		defer cancel()

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"err":"error occured while fetching the food item"})
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc{
	return func(c *gin.Context){

	}
}

func round(num float64) int{

}

func toFixed(num float64, precision int) float64{
	
}

func UpdateFood()gin.HandlerFunc{
	return func(c *gin.Context){

	}
}