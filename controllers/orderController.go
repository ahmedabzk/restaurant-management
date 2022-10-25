package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ahmedabzk/restaurant-management/database"
	"github.com/ahmedabzk/restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)
var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		}

		var allOrders []bson.D

		if err = result.All(ctx, &allOrders); err != nil{
			log.Fatal(err)
		}

		defer cancel()
		c.JSON(http.StatusOK, allOrders)
	}
}

func GetOrder() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		orderId := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id":orderId}).Decode(&order)
		defer cancel()

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		}
		c.JSON(http.StatusOK, order)
		
		

	}
}

func CreateOrders() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func UpdateOrder() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		
	}
}