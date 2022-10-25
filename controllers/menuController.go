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
)
var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.D{})

		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		var allMenus []bson.D

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
			return 
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return 
		}

		err := menuCollection.FindOne(ctx, bson.M{"menu_id":menu.Menu_id}).Decode(&menu)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		}

		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now()).Format(time.RFC3339)
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now()).Format(time.RFC3339)
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()
		menu.Start_date = &time.Time{}
		menu.End_date = &time.Time{}
		
		result, insertErr := menuCollection.InsertOne(ctx, menu)

		defer cancel()
		if insertErr != nil{
			msg := fmt.Sprintf("menu not created")
			c.JSON(http.StatusBadRequest, gin.H{"error":msg})
			return 
		}
		c.JSON(http.StatusOK, result)
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
