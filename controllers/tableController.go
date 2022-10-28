package controllers

import (
	"context"
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

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		result, err := tableCollection.Find(context.TODO(), bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not find the tables in the database"})
			return
		}
		defer cancel()

		var alltables []bson.M

		err = result.All(ctx, &alltables)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, alltables)
	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		tableId := c.Param("table_id")

		var table models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can not find such table in the database"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, table)
	}
}

func CreateTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var table models.Table

		if err := c.BindJSON(&table); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return 
		}
		 
		validationErr := validate.Struct(table)

		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":validationErr.Error()})
			return 
		}
		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		result, insertionErr := tableCollection.InsertOne(ctx, table)
		if insertionErr != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":insertionErr.Error()})
			return 
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		tableId := c.Param("table_id")
		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"table_id": tableId}

		var updateObj primitive.D

		upset := true

		opt := options.UpdateOptions{
			Upsert: &upset,
		}

		if table.Number_of_guests != nil{
			updateObj = append(updateObj, bson.E{"number_of_guests",table.Number_of_guests})
		}

		if table.Table_number != nil{
			updateObj = append(updateObj, bson.E{"table_number",table.Table_number})
		}

		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", table.Updated_at})

		result, err := tableCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set",updateObj},
			},
			&opt,
		)

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return 
		}
		defer cancel()

		c.JSON(http.StatusOK, result)

	}
}
