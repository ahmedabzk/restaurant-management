package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ahmedabzk/restaurant-management/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("STUNNA")

func GenerateAllToken(email, firstName, lastName, uid string) (signedToken, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email: email,
		First_name: firstName,
		Last_name: lastName,
		Uid: uid,
		StandardClaims : jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil{
		log.Panic(err)
		return 
	}
	return token, refreshToken, err
}

func UpdateAllToken(signedToken, signedRefreshToken, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token",Value: signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at",Value: Updated_at})

	filter := bson.M{"user_id":userId}

	upset := true

	opt := options.UpdateOptions{
		Upsert: &upset,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set",Value: updateObj},
		},
		&opt,
	)
	defer cancel()
	if err != nil{
		log.Panic(err)
		return 
	}
	return 
}

func ValidateToken(signedToken string)(claims *SignedDetails, msg string){
	token, err := jwt.ParseWithClaims(
		signedToken,
		SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(signedToken), nil 
		},
	)
	// check if token is invalid
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("token is invalid")
		msg = err.Error()
		return 
	}

	// check if token expired
	if claims.ExpiresAt < time.Now().Local().Unix(){
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return 
	}
	return claims, msg
}
