package data

import (
	"context"
	"fmt"
	"time"

	"github.com/merumaru/marumaru-backend/cfg"
	"github.com/merumaru/marumaru-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllProducts(client *mongo.Client) (*[]models.Product, error) {
	var results []models.Product
	collection := client.Database(cfg.DatabaseName).Collection(cfg.ProductCollection)
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

func GetProductByUserID(client *mongo.Client, id string) (*[]models.Product, error) {
	var results []models.Product
	collection := client.Database("test").Collection("products")
	filter := bson.D{{"sellername", id}}

	cur, err := collection.Find(context.TODO(), filter)
	for cur.Next(context.TODO()) {
		var tmp models.Product
		err := cur.Decode(&tmp)
		if err == nil {
			results = append(results, tmp)
		}
	}
	return &results, err
}

func GetOrderByID(client *mongo.Client, id string) (*models.Order, error) {
	var result models.Order
	collection := client.Database().Collection()
	objID, _ := primitive.ObjectIDFromHex(id) // id is something like "5dc4c0b433f5f1b10da0c599"
	filter := bson.D{{"_id", objID}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Printf("Found a single document: %+v\n", result)
	return &result, err
}

func GetOrderByUserID(client *mongo.Client, id string) (*[]models.Order, error) {
	var results []models.Order
	collection := client.Database("test").Collection("orders")

	filter := bson.M{"$or": []bson.D{bson.D{{"sellername", id}}, bson.D{{"buyername", id}}}}

	cur, err := collection.Find(context.TODO(), filter)
	for cur.Next(context.TODO()) {
		var tmp models.Order
		err := cur.Decode(&tmp)
		if err == nil {
			results = append(results, tmp)
		}
	}
	return &results, err
}

func AddOrder(client *mongo.Client, order *models.Order) error {
	collection := client.Database("test").Collection("orders")
	res, err := collection.InsertOne(context.TODO(), *order)
	fmt.Println("%T", res.InsertedID)
	return err
}

func AddProduct(client *mongo.Client, product *models.Product) error {
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

	id, _ := primitive.ObjectIDFromHex(productID)
	order := models.Order{
		SellerName:   result.SellerID,
		BuyerName:    buyerName,
		ProductID:    id,
		TimeDuration: models.TimeDuration{Start: startDate, End: endDate},
		IsCancelled:  false,
	}
	order.ID = primitive.NewObjectID()
	collectionOrder := client.Database("test").Collection("orders")
	res, err := collectionOrder.InsertOne(context.TODO(), order)
	fmt.Println("%T", res.InsertedID)
	return err
}
