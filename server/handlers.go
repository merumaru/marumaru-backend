package server

import (
	"encoding/json"

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
	ret, _ := json.Marshal(results)
	c.JSON(200, ret)
}

func getProductByIDHandler(c *gin.Context, client *mongo.Client) {
	id := c.Param("id")
	result := data.GetProductByID(client, string(id))
	ret, _ := json.Marshal(result)
	c.JSON(200, ret)

}
