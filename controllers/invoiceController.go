package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ahmedabzk/restaurant-management/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceViewMethod struct {
	Invoice_id       string
	Payment_method   string
	Order_id         string
	Payment_status   *string
	Payment_due      interface{}
	Table_number     interface{}
	Payment_due_date time.Time
	Order_details    interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		result, err := invoiceCollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var allInvoice []bson.M

		if err = result.All(ctx, &allInvoice); err != nil {
			log.Fatal(err)
		}
		defer cancel()
		c.JSON(http.StatusOK, allInvoice)
	}
}

func GetInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func CreateInvoices() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
