package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/merumaru/marumaru-backend/data"
	"github.com/merumaru/marumaru-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
}

func getAllProductsHandler(c *gin.Context, client *mongo.Client) {
	results, err := data.GetAllProducts(client)
	if err != nil {
		c.String(500, "Get Products failed.")
		return
	}
	c.JSON(200, results)
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

func getOrderByIDHandler(c *gin.Context, client *mongo.Client) {
	id := c.Param("id")
	result, err := data.GetOrderByID(client, string(id))
	fmt.Println(result)
	if err != nil {
		c.String(500, "Get Product by ID failed.")
		return
	}
	c.JSON(200, result)
}

func addProductHandler(c *gin.Context, client *mongo.Client) {
	// claims, err := checkLogin(c)
	err := checkLogin_(c, client) // TODO:
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "info": ""})
		return
	}
	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "info": ""})
		return
	}
	// automatically assign an ID
	product.ID = primitive.NewObjectID()
	// product.SellerID = claims.Username // TODO:
	err = data.AddProduct(client, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create product", "info": ""})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Product created!", "info": product.ID})
}

func addOrderHandler(c *gin.Context, client *mongo.Client) {
	// claims, err := checkLogin(c) // TODO: use this line
	err := checkLogin_(c, client) // TODO:

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "info": ""})
		return
	}
	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "info": ""})
		return
	}
	order.ID = primitive.NewObjectID()
	// order.BuyerName = claims.Username // TODO:
	err = data.AddOrder(client, &order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create product", "info": ""})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Order created!", "info": order.ID})
}

func getProductByUserIDHandler(c *gin.Context, client *mongo.Client) {
	id := c.Param("id")
	results, err := data.GetProductByUserID(client, id)
	if err != nil {
		c.String(500, "Get Products failed.")
		return
	}
	c.JSON(200, results)
}

func getOrderByUserIDHandler(c *gin.Context, client *mongo.Client) {
	id := c.Param("id")
	results, err := data.GetOrderByUserID(client, id)
	if err != nil {
		c.String(500, "Get Products failed.")
		return
	}
	c.JSON(200, results)
}

func getOrderByProductIDHandler(c *gin.Context, client *mongo.Client) {
	id := c.Param("id")
	results, err := data.GetOrderByProductID(client, id)
	if err != nil {
		c.String(500, "Get Products failed.")
		return
	}
	c.JSON(200, results)
}

func rentProductHandler(c *gin.Context, client *mongo.Client) {
	claims, err := checkLogin(c)
	if err != nil {
		c.String(500, "Insertion failed.")
		return
	}

	var product models.Product
	id := c.Param("id")
	if err := c.BindJSON(&product); err != nil {
		c.String(400, err.Error())
		return
	}
	buyerName := claims.Username

	dateFormat := "2006-01-02"
	startDate, _ := time.Parse(dateFormat, c.Query("startDate"))
	endDate, _ := time.Parse(dateFormat, c.Query("endDate"))

	err = data.RentProduct(client, string(id), buyerName, startDate, endDate)
	if err != nil {
		c.String(500, "Rent failed.")
		return
	}
	c.String(200, "Product rented.")
}

func editProductHandler(c *gin.Context, client *mongo.Client) {
	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		c.String(400, err.Error())
		return
	}

	id := c.Param("id")

	err := data.Update(client, &product, string(id))
	if err != nil {
		c.String(500, "Update failed.")
		return
	}
	c.String(200, "Product updated.")
}

func cancelProductHandler(c *gin.Context, client *mongo.Client) {

	id := c.Param("id")
	claims, _ := checkLogin(c)
	userID := claims.Username

	if data.CancelOrder(client, string(userID), string(id), false) != nil ||
		data.CancelOrder(client, string(userID), string(id), true) != nil {
		c.String(500, "Cancelation failed.")
		return
	}

	c.String(200, "Product removed.")
}

func getUserByIDHandler(c *gin.Context, client *mongo.Client) {
	id := c.Param("id")
	result, err := data.GetUserByID(client, string(id))
	fmt.Println(result)
	if err != nil {
		c.String(500, "Get User by ID failed. "+err.Error())
		return
	}
	c.JSON(200, result)
}
