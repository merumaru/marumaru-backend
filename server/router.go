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
	router.POST("/user/login", attachDB(client, Signin))
	router.GET("/user/welcome", Welcome)
	router.POST("/user/refresh", Refresh)
	router.POST("/user/signup", attachDB(client, SignUp))
	router.GET("/user", attachDB(client, GetUserByCookie))
	router.POST("/add", attachDB(client, insertProductHandler))
	router.POST("/products/:id/rent", attachDB(client, rentProductHandler))
}

func attachDB(client *mongo.Client, fn func(*gin.Context, *mongo.Client)) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(c, client)
	}
}
