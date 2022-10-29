package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ahmedabzk/restaurant-management/database"
	"github.com/ahmedabzk/restaurant-management/helpers"
	"github.com/ahmedabzk/restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil && recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))

		if err != nil && page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		projectStage := bson.D{
			{
				Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "total_count", Value: 1},
					{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
				}}}
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, projectStage})
			defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while trying to list user items"})
			return 
		}

		var allUsers []bson.M

		if err := result.All(ctx, &allUsers); err != nil{
		log.Fatal(err)
		}
		defer cancel()

		c.JSON(http.StatusOK, allUsers)
	}

}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		userId := c.Param("user_id")

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, user)
	}
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		}

		// check if email already exists
		count, err := userCollection.CountDocuments(ctx, bson.M{"email":user.Email})
		defer cancel()

		if err != nil{
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return 
		}

		// hash password
		password := HashPassword(*user.Password)
		user.Password = &password

		// check if phone number exists
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone":user.Phone})
		defer cancel()
		if err != nil{
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return 
		}
		if count > 0{
			c.JSON(http.StatusBadRequest, gin.H{"error":"email or phone number already exists"})
		}
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// generate token
		token, refreshToken, _ := helpers.GenerateAllToken(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		// insert user to the database

		result, insertionErr := userCollection.InsertOne(ctx, user)
		defer cancel()

		if insertionErr != nil{
			log.Fatal(insertionErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error":insertionErr.Error()})
			return 
		}

		c.JSON(http.StatusOK, result)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		var foundUser models.User

		if err := c.BindJSON(&user); err != nil{
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		}

		err := userCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":"can not login, check your email or password again"})
			return 
		}

		// verify password
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		if passwordIsValid != true{
			c.JSON(http.StatusBadRequest, gin.H{"error":msg})
			return 
		}
		// generate token
		token, refreshToken, _ := helpers.GenerateAllToken(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)
		// update token

		helpers.UpdateAllToken(token, refreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, foundUser)
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil{
		log.Fatal(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))

	check := true
	msg := ""

	if err != nil{
		msg = fmt.Sprintf("email or password is wrong")
		check = false
	}

	return check, msg
}
