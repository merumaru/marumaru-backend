package main

import "github.com/codechrysalis/go.pokemon-api/server"

func main() {
	router := server.CreateRouter()
	server.StartServer(router)
}
