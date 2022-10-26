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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InvoiceViewFormat struct {
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
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		invoiceId := c.Param("invoice_id")
		var invoice models.Invoice

		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id":invoiceId}).Decode(&invoice)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while trying to fetch the item"})
			return 
		}

		var invoiceView InvoiceViewFormat

		allOrderItems, err := ItemsByOrder(invoice.Order_id)
		invoiceView.Order_id = invoice.Order_id
		invoiceView.Payment_due_date = invoice.Payment_due_date

		invoiceView.Payment_method = "null"
		if invoice.Payment_method != nil{
			invoiceView.Payment_method = *invoice.Payment_method
		}
		invoiceView.Invoice_id = invoice.Invoice_id
		invoiceView.Payment_status = *&invoice.Payment_status
		invoiceView.Payment_due = allOrderItems[0]["payment_due"]
		invoiceView.Table_number = allOrderItems[0]["table_number"]
		invoiceView.Order_details = allOrderItems[0]["order_items"]

		c.JSON(http.StatusOK, invoiceView)
	}
}

func CreateInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var invoice models.Invoice
		if err := c.BindJSON(&invoice); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return 
		}
		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id":invoice.Order_id}).Decode(&order)
		defer cancel()

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"did not find the order"})
			return 
		}
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(),100*time.Second)

		invoiceId := c.Param("invoice_id")

		var invoice models.Invoice

		if err := c.BindJSON(&invoice); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return
		}

		filter := bson.M{"invoice_id":invoiceId}

		var updateObj primitive.D

		if invoice.Payment_method != nil{

		}
		if invoice.Payment_status != nil{

		}

		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", invoice.Updated_at})

		invoice.ID = primitive.NewObjectID()
		invoice.Invoice_id = invoice.ID.Hex()

		upset := true

		opt := options.UpdateOptions{
			Upsert: &upset,
		}

		status := "PENDING"
		if invoice.Payment_status == nil{
			invoice.Payment_status = &status
		}

		result, err := invoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"&set", updateObj},
			},
			&opt,
		)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"could not update the invoice"})
			return 
		}
		c.JSON(http.StatusOK, result)
	}
}
