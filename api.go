package main

import "github.com/merumaru/marumaru-backend/server"

func main() {
	router := server.CreateRouter()
	server.StartServer(router)

}
