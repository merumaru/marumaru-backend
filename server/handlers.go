package server

import (
	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	c.String(200, "This is indeed a landing page")
}

func loginPage(c *gin.Context) {
	c.String(200, "loginPage")
}

func listPage(c *gin.Context) {
	c.String(200, "listPage")
}

func getAllProducts(c *gin.Context) {
	c.String(200, "getAllProducts")
}

func getProductByID(c *gin.Context) {
	id := c.Param("id")
	c.String(200, id)
}
