package main

import (
	"context"
	"fmt"
	"log"

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

	ash := Trainer{"Ash", 10, "Pallet Town"}
	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

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
