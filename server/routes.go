package server

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupRoutes(router *gin.Engine) {
	dbClient := Connect2DB()
	router.GET("/", hello)
	router.GET("/login", loginPage)
	router.GET("/list", attachDB(dbClient, listPage))

	router.GET("/products", attachDB(dbClient, getAllProductsHandler))
	router.GET("/products/:id", attachDB(dbClient, getProductByIDHandler))
	router.POST("/products", attachDB(dbClient, addProductHandler))
	router.POST("/products/:id/rent", attachDB(dbClient, rentProductHandler))
	router.PATCH("/products/:id/edit", attachDB(dbClient, editProductHandler))
	router.PATCH("/products/:id/cancel", attachDB(dbClient, cancelProductHandler))
	router.PATCH("/products/:id/recommendations", attachDB(dbClient, getRecommendationsHandler))

	router.POST("/orders", attachDB(dbClient, addOrderHandler))
	router.GET("/orders/:id", attachDB(dbClient, getOrderByIDHandler))

	router.GET("/users/:id/products", attachDB(dbClient, getProductByUserIDHandler))
	router.GET("/users/:id/orders", attachDB(dbClient, getOrderByUserIDHandler))
	router.GET("/users/:id", attachDB(dbClient, getUserByIDHandler))

	router.POST("/login", attachDB(dbClient, signIn))
	router.GET("/welcome", welcome)
	router.POST("/refresh", refresh)
	router.POST("/signup", attachDB(dbClient, signUp))
	router.GET("/cookie", attachDB(dbClient, getUserByCookie))
}

func attachDB(client *mongo.Client, fn func(*gin.Context, *mongo.Client)) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(c, client)
	}
}
