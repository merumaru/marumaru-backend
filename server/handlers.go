package server

import (
	"fmt"
	"log"
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

// TODO:
func getAllProductsHandler(c *gin.Context, client *mongo.Client) {
	results, err := data.GetAllProducts(client)
	if err != nil {
		log.Println("Get products failed : ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Get Products failed : " + err.Error(), "info": ""})
		return
	}
	log.Println("Fetched results successfully")
	c.JSON(200, results)
}

// TODO:
func getProductByIDHandler(c *gin.Context, client *mongo.Client) {
	id := c.Param("id")
	result, err := data.GetProductByID(client, string(id))
	fmt.Println(result)
	if err != nil {
		log.Println(err)
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

// TODO:
func addProductHandler(c *gin.Context, client *mongo.Client) {
	userCookie := getUserByCookie(c)
	if userCookie == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User is not signed in", "info": ""})
		return
	}

	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "info": ""})
		return
	}
	// automatically assign an ID
	product.ID = primitive.NewObjectID()
	product.SellerID = userCookie.Username

	err := data.AddProduct(client, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create product", "info": ""})
		return
	}
	log.Println("Added product successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "Product created!", "info": product.ID})
}

func addOrderHandler(c *gin.Context, client *mongo.Client) {
	userCookie := getUserByCookie(c)
	if userCookie == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User is not signed in", "info": ""})
		return
	}

	var order models.Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "info": ""})
		return
	}
	order.ID = primitive.NewObjectID()
	order.BuyerID = userCookie.Username // TODO:
	if order.BuyerID == order.SellerID {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot buy your own product!!", "info": ""})
		return
	}
	err := data.AddOrder(client, &order)
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

// TODO:
func rentProductHandler(c *gin.Context, client *mongo.Client) {
	userCookie := getUserByCookie(c)
	if userCookie == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User is not signed in", "info": ""})
		return
	}

	id := c.Param("id")
	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		c.String(400, err.Error())
		return
	}
	buyerName := userCookie.Username

	dateFormat := "2006-01-02"
	startDate, _ := time.Parse(dateFormat, c.Query("startDate"))
	endDate, _ := time.Parse(dateFormat, c.Query("endDate"))

	if err := data.RentProduct(client, string(id), buyerName, startDate, endDate); err != nil {
		c.String(500, "Rent failed.")
		return
	}
	c.String(200, "Product rented.")
}

// TODO:
func editProductHandler(c *gin.Context, client *mongo.Client) {
	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		c.String(400, err.Error())
		return
	}

	id := c.Param("id")

	err := data.Update(client, &product, string(id))
	if err != nil {
		log.Println(err)
		c.String(500, "Update failed.")
		return
	}
	c.String(200, "Product updated.")
}

func cancelProductHandler(c *gin.Context, client *mongo.Client) {

	id := c.Param("id")

	userCookie := getUserByCookie(c)
	if userCookie == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User is not signed in", "info": ""})
		return
	}
	userID := userCookie.Username

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
		c.String(500, "Get User by ID failed."+err.Error())
		return
	}
	c.JSON(200, result)
}
