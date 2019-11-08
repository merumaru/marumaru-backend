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
	router.POST("/login", attachDB(client, Signin))
	router.GET("/welcome", Welcome)
	router.POST("/refresh", Refresh)
	router.POST("/signup", attachDB(client, SignUp))

}

func attachDB(client *mongo.Client, fn func(*gin.Context, *mongo.Client)) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(c, client)
	}
}
