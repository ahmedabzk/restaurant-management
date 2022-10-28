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

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

func GetUsers() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}

}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		userId := c.Param("user_id")

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id":userId}).Decode(&user)

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return 
		}
		defer cancel()
		c.JSON(http.StatusOK, user)
	}
}

func Signup() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func Login() gin.HandlerFunc{
	return func(ctx *gin.Context) {

	}
}

func HashPassword(password string) string{

}

func VerifyPassword(userPassword, providedPassword string) (bool, string){
	
}