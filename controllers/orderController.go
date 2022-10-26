package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ahmedabzk/restaurant-management/database"
	"github.com/ahmedabzk/restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		var allOrders []bson.D

		if err = result.All(ctx, &allOrders); err != nil {
			log.Fatal(err)
		}

		defer cancel()
		c.JSON(http.StatusOK, allOrders)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		orderId := c.Param("order_id")
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, order)

	}
}

func CreateOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var order models.Order
		var table models.Table

		if err := c.BindJSON(&order); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return 
		}

		if order.Table_id != nil{
			err := tableCollection.FindOne(ctx, bson.M{"table_id":order.Table_id}).Decode(&table)
			defer cancel()
			if err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
				return 
			}
		}
		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()

		result, insertErr := orderCollection.InsertOne(ctx, order)

		if insertErr != nil{
			msg := fmt.Sprintf("could not create the new order")
			c.JSON(http.StatusBadRequest, gin.H{"error":msg})
			return 
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var table models.Table
		var order models.Order

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var updateObj primitive.D

		orderId := c.Param("order_id")

		if err := c.BindJSON(&order); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return 
		}
		if order.Table_id != nil{
			err := menuCollection.FindOne(ctx, bson.M{"table_id":order.Table_id}).Decode(&table)

			defer cancel()

			if err != nil{
				c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
				return 
			}
			updateObj = append(updateObj, bson.E{"menu", order.Table_id})
		}
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", order.Updated_at})

		upset := true
		filter := bson.M{"order_id":orderId}

		opt := options.UpdateOptions{
			Upsert: &upset,
		}

		result, err := orderCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"&set", updateObj},
			},
			&opt,
		)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
