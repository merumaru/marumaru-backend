package server

import (
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/merumaru/marumaru-backend/data"
	"github.com/merumaru/marumaru-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	jwtUserIdKey   = "id"
	jwtUserNameKey = "username"
	authMiddleware *jwt.GinJWTMiddleware
	mongoDBClient  *mongo.Client //initialized in router when it is being created
)

func setUpJWT() {
	// the jwt middleware
	var err error
	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:          "MaruMaru API",
		Key:            []byte("secret key"),
		IdentityKey:    jwtUserNameKey,
		Timeout:        24 * 3 * time.Hour,
		MaxRefresh:     24 * 3 * time.Hour,
		SendCookie:     true,
		SecureCookie:   false,
		CookieHTTPOnly: true,
		CookieName:     "token",
		TokenLookup:    "cookie:token",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					jwtUserIdKey:   v.ID.Hex(),
					jwtUserNameKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			objID, err := primitive.ObjectIDFromHex(claims[jwtUserIdKey].(string))
			if err != nil {
				return nil
			}
			return &models.User{
				ID:       objID,
				Username: claims[jwtUserNameKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			log.Println(c.Request.URL)
			loginVals := new(struct {
				Username string `form:"username" json:"username" binding:"required"`
				Password string `form:"password" json:"password" binding:"required"`
			})
			if err := c.ShouldBind(loginVals); err != nil {
				log.Println("missing values when loginng in")
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if user, err := data.GetUserByUserName(mongoDBClient, userID); err != nil {
				log.Printf("Error in getting user details for '%s' from DB: %s", userID, err.Error())
				return "", jwt.ErrFailedAuthentication
			} else {
				if user == nil {
					log.Printf("User object is nil for user '%s'", userID)
					return "", jwt.ErrFailedAuthentication
				}
				if user.Password != password {
					log.Printf("Incorrect password for user '%s'", userID)
					return "", jwt.ErrFailedAuthentication
				}
				log.Printf("Successfully logged in user '%s'", userID)
				return user, nil
			}
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	})
	if err != nil {
		panic("failed to create jwt middleware")
	}
	log.Println("Successfully set up JWT")
}

func setupRoutes(router *gin.Engine, databaseURL, databaseName string) {
	dbClient := Connect2DB(databaseURL)
	mongoDBClient = dbClient //global variable to save mongo client. used to access DB in JWT.

	setUpJWT()
	router.POST("/users/login", authMiddleware.LoginHandler)
	router.POST("/users/signup", attachDB(dbClient, signUp))

	auth := router.Group("")
	auth.Use(authMiddleware.MiddlewareFunc())

	auth.GET("/", hello)
	auth.GET("/login", loginPage)                   // TODO: not used
	auth.GET("/list", attachDB(dbClient, listPage)) // TODO: not used

	auth.GET("/products", attachDB(dbClient, getAllProductsHandler))
	auth.GET("/products/:id", attachDB(dbClient, getProductByIDHandler))
	auth.POST("/products", attachDB(dbClient, addProductHandler))
	auth.POST("/products/:id/rent", attachDB(dbClient, rentProductHandler))
	auth.PATCH("/products/:id/edit", attachDB(dbClient, editProductHandler))
	auth.PATCH("/products/:id/cancel", attachDB(dbClient, cancelProductHandler))

	auth.POST("/orders", attachDB(dbClient, addOrderHandler))
	auth.GET("/orders/:id", attachDB(dbClient, getOrderByIDHandler))

	auth.GET("/users/:id/products", attachDB(dbClient, getProductByUserIDHandler))
	auth.GET("/users/:id/orders", attachDB(dbClient, getOrderByUserIDHandler))
	auth.GET("/users/:id", attachDB(dbClient, getUserByIDHandler))
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	//TODO remove this
	// router.GET("/usercookie", attachDB(dbClient, getUserByCookie))
}

func attachDB(client *mongo.Client, fn func(*gin.Context, *mongo.Client)) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(c, client)
	}
}
