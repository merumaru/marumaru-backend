package data

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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
	collection := client.Database(cfg.DatabaseName).Collection(cfg.ProductCollection)
	objID, _ := primitive.ObjectIDFromHex(id) // id is something like "5dc4c0b433f5f1b10da0c599"
	filter := bson.D{{"_id", objID}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Printf("Found a single document: %+v\n", result)
	return &result, err
}

func GetProductByUserID(client *mongo.Client, id string) (*[]models.Product, error) {
	var results []models.Product
	collection := client.Database(cfg.DatabaseName).Collection(cfg.ProductCollection)
	filter := bson.D{{"sellerid", id}}

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
	collection := client.Database(cfg.DatabaseName).Collection(cfg.OrderCollection)
	objID, _ := primitive.ObjectIDFromHex(id) // id is something like "5dc4c0b433f5f1b10da0c599"
	filter := bson.D{{"_id", objID}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Printf("Found a single document: %+v\n", result)
	return &result, err
}

func GetOrderByProductID(client *mongo.Client, id string) (*[]models.Order, error) {
	var results []models.Order
	collection := client.Database(cfg.DatabaseName).Collection(cfg.OrderCollection)
	objID, _ := primitive.ObjectIDFromHex(id) // id is something like "5dc4c0b433f5f1b10da0c599"
	filter := bson.D{{"productid", objID}}

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

func GetOrderByUserID(client *mongo.Client, id string) (*[]models.Order, error) {
	var results []models.Order
	collection := client.Database(cfg.DatabaseName).Collection(cfg.OrderCollection)

	filter := bson.D{{"buyerid", id}}
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
	collection := client.Database(cfg.DatabaseName).Collection(cfg.OrderCollection)
	res, err := collection.InsertOne(context.TODO(), *order)
	fmt.Println("%T", res.InsertedID)
	return err
}

func AddProduct(client *mongo.Client, product *models.Product) error {
	collection := client.Database(cfg.DatabaseName).Collection(cfg.ProductCollection)
	res, err := collection.InsertOne(context.TODO(), *product)
	fmt.Println("%T", res.InsertedID)
	return err
}

func RentProduct(client *mongo.Client, productID string, buyerName string, startDate time.Time, endDate time.Time) error {
	var result models.Product
	collection := client.Database(cfg.DatabaseName).Collection(cfg.ProductCollection)
	objID, _ := primitive.ObjectIDFromHex(productID)
	filter := bson.D{{"_id", objID}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	// id, _ := primitive.ObjectIDFromHex(productID)
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
	collection := client.Database(cfg.DatabaseName).Collection(cfg.ProductCollection)
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": bson.M{"$eq": objID}}
	_, err := collection.UpdateOne(context.TODO(), filter, *product)
	return err
}

func CancelOrder(client *mongo.Client, userID string, id string, whoCancelled bool) error {
	// 0 --> Buyer, 1--> Seller cancelled
	t := time.Now()
	collection := client.Database("test").Collection("orders")
	usrID, _ := primitive.ObjectIDFromHex(userID)

	objID, _ := primitive.ObjectIDFromHex(id)

	var filter bson.M
	if whoCancelled == false {
		filter = bson.M{"BuyerID": bson.M{"$eq": usrID}, "ProductID": bson.M{"$eq": objID}, "TimeDuration.Start": bson.M{"$gte": t}}
	} else {
		filter = bson.M{"SellerID": bson.M{"$eq": usrID}, "ProductID": bson.M{"$eq": objID}, "TimeDuration.Start": bson.M{"$gte": t}}
	}
	update := bson.M{"$set": bson.M{"IsCancelled": true}}
	_, err := collection.UpdateMany(context.TODO(), filter, update)
	return err
}

func addProductToRecSysDB(productID string, imageURL string) error {
	requestBody, _ := json.Marshal(map[string]string{
		"url": imageURL,
	})
	resp, err := http.Post(fmt.Sprintf("http://34.83.27.35:5000/%s/addImage", productID), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return err
}

func AddUser(client *mongo.Client, user *models.User) error {

	u, err := GetUserByUserName(client, user.Username)
	if err != nil {
		log.Println("Error occured in getting user by username : ", err.Error())
		return err
	} else if u != nil {
		return errors.New("User with same Username exists")
	}

	collection := client.Database(cfg.DatabaseName).Collection(cfg.UserCollection)
	user.ID = primitive.NewObjectID()
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println("Error occured in inserting new user : ", err.Error())
		return err
	}
	fmt.Println("Inserted new user with ID: ", insertResult.InsertedID)
	return nil
}

func GetUserByUserName(client *mongo.Client, username string) (*models.User, error) {
	var result models.User
	collection := client.Database(cfg.DatabaseName).Collection(cfg.UserCollection)
	filter := bson.M{"username": username}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	fmt.Printf("Found a single document: %+v\n", result)
	return &result, nil
}

func GetUserByID(client *mongo.Client, id string) (*models.User, error) {
	var result models.User
	collection := client.Database(cfg.DatabaseName).Collection(cfg.UserCollection)
	objID, _ := primitive.ObjectIDFromHex(id) // id is something like "5dc4c0b433f5f1b10da0c599"
	filter := bson.D{{"_id", objID}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	fmt.Printf("Found a single document: %+v\n", result)
	return &result, nil
}

func GetRecommendations(client *mongo.Client, productID string) (*[]models.Product, error) {
	var results []models.Product
	recommendation := new(models.Recommendation)

	resp, err := http.Get(fmt.Sprintf("http://34.83.27.35:5000/%s/similarProducts", productID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(recommendation)

	for _, id := range recommendation.ProductList {
		product, err := GetProductByID(client, id)
		if err != nil {
			return nil, err
		}
		results = append(results, *product)
	}
	return &results, err
}
