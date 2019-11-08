package data

import (
	"context"
	"fmt"
	"log"

	"github.com/merumaru/marumaru-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllProducts(client *mongo.Client) *[]models.Product {
	var results []models.Product
	collection := client.Database("test").Collection("product")
	cur, _ := collection.Find(context.TODO(), bson.D{})
	for cur.Next(context.TODO()) {
		var tmp models.Product
		err := cur.Decode(&tmp)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		results = append(results, tmp)
	}
	return &results
}

func GetProductByID(client *mongo.Client, id string) *models.Product {
	var result models.Product
	collection := client.Database("test").Collection("product")
	objID, _ := primitive.ObjectIDFromHex(id) // id is something like "5dc4c0b433f5f1b10da0c599"
	filter := bson.D{{"_id", objID}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Printf("Found a single document: %+v\n", result)
	if err != nil {
		log.Fatal(err)
	}
	return &result
}
