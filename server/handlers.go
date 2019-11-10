package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/merumaru/marumaru-backend/data"
	"go.mongodb.org/mongo-driver/mongo"
)

func hello(c *gin.Context) {
	c.String(200, "This is indeed a landing page")
}

func loginPage(c *gin.Context) {
	c.String(200, "loginPage")
}

func listPage(c *gin.Context, client *mongo.Client) {
	c.String(200, "listPage")
	// client.Database("test").Collection("product")
}

func getAllProductsHandler(c *gin.Context, client *mongo.Client) {
	results := data.GetAllProducts(client)
	c.JSON(200, *results)
}

func getProductByIDHandler(c *gin.Context, client *mongo.Client) {
	id := c.Param("id")
	result, err := data.GetProductByID(client, string(id))
	fmt.Println(result)
	if err != nil {
		c.String(500, "Get Product by ID failed.")
		return
	}
	c.JSON(200, result)

}
