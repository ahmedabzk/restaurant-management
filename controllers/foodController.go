package controllers

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/ahmedabzk/restaurant-management/database"
	"github.com/ahmedabzk/restaurant-management/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
// var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var validate = validator.New()

func GetFoods() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1{
			recordPerPage = 10
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1{
			page = 1
		}

		startIndex := (page - 1)* recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		// mongodb aggregation

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}}, {Key: "total_count", Value: bson.D{{Key: "$sum", Value: "1"}}}, {Key: "data", Value: bson.D{{Key: "$push",Value: "$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "total_count", Value: 1},
					{Key: "food_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
				}}}
		

		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,groupStage, projectStage})
		defer cancel()

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		var allFoods []bson.M

		if err := result.All(ctx, &allFoods); err != nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allFoods[0])
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
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var food models.Food
		var menu models.Menu

		if err := c.BindJSON(&food); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return 
		}

		validationErr := validate.Struct(food)

		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return 
		}
		err := menuCollection.FindOne(ctx, bson.M{"menu_id":food.Menu_id}).Decode(&menu)
		defer cancel()

		if err != nil{
			msg := fmt.Sprintf("menu not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num

		result, insertErr := foodCollection.InsertOne(ctx, food)
		defer cancel()
		if insertErr != nil{
			msg := fmt.Sprintf("food item not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return 
		}
		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int{
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64{
	output := math.Pow(10, float64(precision))
	return float64(round(num*output))/output
}

func UpdateFood()gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var menu models.Menu
		var food models.Food

		var foodID = c.Param("food_id")

		if err := c.BindJSON(&food); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		}

		var updateObj primitive.D

		if food.Name != nil{
			updateObj = append(updateObj, bson.E{"name", food.Name})
		}
		if food.Price != nil{
			updateObj = append(updateObj, bson.E{"price", food.Price})
		}

		if food.Food_image != nil{
			updateObj = append(updateObj, bson.E{"food_image", food.Food_image})
		}

		if food.Menu_id != nil{
			err := menuCollection.FindOne(ctx, bson.M{"menu_id":food.Menu_id}).Decode(&menu)
			if err != nil{
				c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			}
			defer cancel()

			updateObj = append(updateObj, bson.E{"menu", food.Price})
		}
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", food.Updated_at})

		upsert := true
		filter := bson.M{"food_id":foodID}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := foodCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}