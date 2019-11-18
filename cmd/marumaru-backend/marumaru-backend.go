package main

import (
	"flag"
	"fmt"

	"github.com/merumaru/marumaru-backend/server"
)

func main() {
	dbURL := flag.String("mongodb-url", "mongodb://localhost:27017", "URL to connet to mongodb database")
	dbName := flag.String("database", "testing", "Database name in mongodb")
	flag.Parse()
	fmt.Println("Database url  ", *dbURL)
	router := server.CreateRouter(*dbURL, *dbName)
	server.StartServer(router)
}
