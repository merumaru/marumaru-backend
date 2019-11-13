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
	router.GET("/products-user/:id", attachDB(client, getProductByUserIDHandler))
	router.POST("/products", attachDB(client, addProductHandler))
	router.POST("/products/:id/rent", attachDB(client, rentProductHandler))
	router.PATCH("/products/:id/edit", attachDB(client, editProductHandler))
	router.PATCH("/products/:id/cancel", attachDB(client, cancelProductHandler))

	router.POST("/orders", attachDB(client, addOrderHandler))
	router.GET("/orders-user/:id", attachDB(client, getOrderByUserIDHandler))
	router.GET("/orders/:id", attachDB(client, getOrderByIDHandler))

	router.POST("/users/login", attachDB(client, Signin))
	router.GET("/users/welcome", Welcome)
	router.POST("/users/refresh", Refresh)
	router.POST("/users/signup", attachDB(client, SignUp))
	router.GET("/users/user/:id", attachDB(client, getUserByIDHandler))
	router.GET("/users", attachDB(client, GetUserByCookie))
}

func attachDB(client *mongo.Client, fn func(*gin.Context, *mongo.Client)) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(c, client)
	}
}
