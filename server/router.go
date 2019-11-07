package server

import (
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine) {
	// It is good practice to version your API from the start
	v1 := router.Group("/api/v1")

	v1.GET("/", loginPage)
	v1.GET("/mainpage", mainPage)
	v1.GET("/products", getAllProducts)
	v1.GET("/products/:id", getProductByID)

}
