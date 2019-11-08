package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/merumaru/marumaru-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	collection := client.Database("test").Collection("trainers")
	formate := "2006-01-02"
	start, _ := time.Parse(formate, "2018-05-31")
	end, _ := time.Parse(formate, "2018-05-31")
	product1 := models.Product{[]string{"url1", "url2"}, "name", "desp",
		100, models.TimeDuration{start, end}, []int{models.Book}}
	product2 := models.Product{[]string{"url1", "url2"}, "name", "desp",
		100, models.TimeDuration{start, end}, []int{models.Book}}
	products := []interface{}{product1, product2}
	// insert one
	// insertResult, err := collection.InsertOne(context.TODO(), product1)
	// insert many
	insertResult, err := collection.InsertMany(context.TODO(), products)
	if err != nil {
		log.Fatal(err)
	}

	// find by id
	fmt.Println("Inserted a single document: ", insertResult.InsertedIDs)
	objID, _ := primitive.ObjectIDFromHex("5dc4c0b433f5f1b10da0c599")
	filter := bson.D{{"_id", objID}}
	var result models.Product
	var results []models.Product

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Printf("Found a single document: %+v\n", result)

	// find all
	cur, _ := collection.Find(context.TODO(), bson.D{})
	for cur.Next(context.TODO()) {
		var tmp models.Product
		err = cur.Decode(&tmp)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		results = append(results, tmp)
	}
	fmt.Println(results)
	if err != nil {
		log.Fatal(err)
	}

}

// var (
// 	attacks models.AttackList
// 	pokemon models.PokemonList
// 	types   models.TypeList
// 	box     = packr.New("assets", "../assets")
// )

// // Pokemon contains all available Pokemon
// func Pokemon() *models.PokemonList {
// 	return &pokemon
// }

// // Attacks contains all available attacks
// func Attacks() *models.AttackList {
// 	return &attacks
// }

// // Types contains all available Types
// func Types() *models.TypeList {
// 	return &types
// }

// func loadFile(name string) []byte {
// 	res, err := box.Find(name)

// 	if err != nil {
// 		panic(err)
// 	}

// 	return res
// }

// func init() {
// 	Reload()
// }

// // Reload the data from the json files
// func Reload() {
// 	json.Unmarshal(loadFile("attacks.json"), &attacks)
// 	json.Unmarshal(loadFile("pokemon.json"), &pokemon)
// 	json.Unmarshal(loadFile("types.json"), &types)
// }
