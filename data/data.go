package data

import (
	"context"
	"fmt"
	"time"

	"github.com/merumaru/marumaru-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllProducts(client *mongo.Client) (*[]models.Product, error) {
	var results []models.Product
	collection := client.Database("test").Collection("products")
	cur, err := collection.Find(context.TODO(), bson.D{})
	for cur.Next(context.TODO()) {
		var tmp models.Product
		err := cur.Decode(&tmp)
		if err == nil {
			results = append(results, tmp)
		}
	}
	return &results, err
}

func GetProductByID(client *mongo.Client, id string) (*models.Product, error) {
	var result models.Product
	collection := client.Database("test").Collection("products")
	objID, _ := primitive.ObjectIDFromHex(id) // id is something like "5dc4c0b433f5f1b10da0c599"
	filter := bson.D{{"_id", objID}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Printf("Found a single document: %+v\n", result)
	return &result, err
}

func Insert(client *mongo.Client, product *models.Product) error {
	collection := client.Database("test").Collection("products")
	res, err := collection.InsertOne(context.TODO(), *product)
	fmt.Println("%T", res.InsertedID)
	return err
}

func RentProduct(client *mongo.Client, productID string, buyerName string, startDate time.Time, endDate time.Time) error {
	var result models.Product
	collection := client.Database("test").Collection("products")
	objID, _ := primitive.ObjectIDFromHex(productID)
	filter := bson.D{{"_id", objID}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	order := models.Order{
		SellerID:     result.SellerID,
		BuyerID:      buyerName,
		ProductID:    productID,
		TimeDuration: models.TimeDuration{Start: startDate, End: endDate},
		IsCancelled:  false,
	}
	order.ID = primitive.NewObjectID()
	collectionOrder := client.Database("test").Collection("orders")
	res, err := collectionOrder.InsertOne(context.TODO(), order)
	fmt.Println("%T", res.InsertedID)
	return err
}

func Update(client *mongo.Client, product *models.Product, id string) error {
	collection := client.Database("test").Collection("products")
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": bson.M{"$eq": objID}}
	_,err := collection.UpdateOne(context.TODO(), filter, *product)
	return err
}