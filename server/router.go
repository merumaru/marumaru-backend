package server

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupRoutes(router *gin.Engine) {
	// v1 := router.Group("/api/v1")
	client := Connect2DB()
	router.GET("/", hello)
	router.GET("/login", loginPage)
	router.GET("/list", attachDB(client, listPage))
	router.GET("/products", attachDB(client, getAllProductsHandler))
	router.GET("/products/:id", attachDB(client, getProductByIDHandler))

}

func attachDB(client *mongo.Client, fn func(*gin.Context, *mongo.Client)) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(c, client)
	}
}
