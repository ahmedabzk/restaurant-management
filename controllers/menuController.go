package controllers

import (
	"context"
	"fmt"
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
var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.M{})

		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		var allMenus []bson.M

		err = result.All(ctx, &allMenus)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		}
		c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		// get the id from the params
		menuId := c.Param("menu_id")

		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id":menuId}).Decode(&menu)
		defer cancel()

		if err != nil{
			msg := fmt.Sprintf("error occured while fetching the menu")
			c.JSON(http.StatusBadRequest, gin.H{"error":msg})
			return 
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			
		}


		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()
		
		
		result, insertErr := menuCollection.InsertOne(ctx, menu)

		defer cancel()
		if insertErr != nil{
			msg := fmt.Sprintf("menu not created")
			c.JSON(http.StatusBadRequest, gin.H{"error":msg})
			return 
		}
		c.JSON(http.StatusOK, result)
		defer cancel()
	}
}
func inTimeSpan(start, end, check time.Time) bool{
	return start.After(time.Now()) && end.After(start)
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		}
		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id":menuId}

		var updateObj primitive.D

		if menu.Start_date != nil && menu.End_date != nil{
			if !inTimeSpan(*menu.Start_date, *menu.End_date, time.Now()){
				msg := fmt.Sprintf("please retype the time")
				c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
				defer cancel()
				return 
			}
			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.Start_date})
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.End_date})

			if menu.Name != nil{
				updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
			}
			if menu.Category != nil{
				updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
			}
			menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.Updated_at})

			upsert := true

			opt := options.UpdateOptions{
				Upsert: &upsert,
			}

			result, err := menuCollection.UpdateOne(
				ctx,
				filter,
				bson.D{
					{Key: "$set", Value: updateObj},
				},
				&opt,
			)

			defer cancel()

			if err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			}
			
			c.JSON(http.StatusOK, result)
		}


	}
}
