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
	router.GET("/product/:id", attachDB(client, getProductByIDHandler))
	router.GET("/product-user/:id", attachDB(client, getProductByUserIDHandler))
	router.GET("/order-user/:id", attachDB(client, GetOrderByUserIDHandler))
	router.GET("/order/:id", attachDB(client, getOrderByIDHandler))
	router.POST("/addproduct", attachDB(client, addProductHandler))
	router.POST("/addorder", attachDB(client, addOrderHandler))

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
