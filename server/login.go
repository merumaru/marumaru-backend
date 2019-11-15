package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/merumaru/marumaru-backend/cfg"
	models "github.com/merumaru/marumaru-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: what is this
// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

// Create handler which needs db client
type HandlerFuncWithDB func(*gin.Context, *mongo.Client)

func checkForUser(filter bson.M, client *mongo.Client, username string) (int, string, models.User) {
	collection := client.Database(cfg.DatabaseName).Collection(cfg.UserCollection)
	result := models.User{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		// No result
		if err == mongo.ErrNoDocuments {
			return http.StatusNotFound, fmt.Sprintf("No such user name %s!", username), result
		} else {
			// Other problem
			return http.StatusInternalServerError, err.Error(), result
		}
	}
	return http.StatusOK, "", result
}

// Signin add a user cookie
func signIn(c *gin.Context, client *mongo.Client) {
	var creds models.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.String(400, err.Error())
		return
	}

	filter := bson.M{"username": creds.Username}
	status, msg, result := checkForUser(filter, client, creds.Username)
	if status != http.StatusOK {
		c.String(status, msg)
		return
	}
	// expectedPassword, ok := users[creds.Username]
	// Check password
	if result.Password != creds.Password {
		c.String(http.StatusUnauthorized, "Password is incorrect!")
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 10 minutes
	expirationTime := time.Now().Add(10 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &models.Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "token",
	// 	Value:   tokenString,
	// 	Expires: expirationTime,
	// })
	c.SetCookie("token", tokenString, expirationTime.Second(), "/", "", false, false)
	c.JSON(http.StatusOK, gin.H{"message": "Set cookie successfully", "info": result.ID.Hex()})
	return
}

// Signup adds a new user
func signUp(c *gin.Context, client *mongo.Client) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.String(400, err.Error())
		return
	}
	// Check input
	if user.Username == "" || user.Password == "" || user.Email == "" {
		c.String(400, "Invalid input")
		return
	}
	collection := client.Database(cfg.DatabaseName).Collection(cfg.UserCollection)
	filter := bson.M{"username": user.Username}
	result := models.User{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err == nil {
		c.String(http.StatusBadRequest, "Username exists!")
		return
	}
	if err != mongo.ErrNoDocuments {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	// expectedPassword, ok := users[creds.Username]
	// Add user
	user.ID = primitive.NewObjectID()
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "Sign up successfully!")
}

// Welcome tests login
func welcome(c *gin.Context) {
	claims, err := checkLogin(c)
	if err != nil {
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("Welcome %s!", claims.Username))
}

// Refresh the token in the background by the client application.
func refresh(c *gin.Context) {
	claims, err := checkLogin(c)
	if err != nil {
		return
	}
	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.String(http.StatusBadRequest, "A new token will only be issued if the old token is within 30 seconds of expiry")
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("token", tokenString, expirationTime.Second(), "/", "", false, false)
	c.String(http.StatusOK, "Refresh a new token!")
	return
}

// GetUserByCookie returns the whole user struct by cookie
func getUserByCookie(c *gin.Context, client *mongo.Client) {
	claims, err := checkLogin(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	collection := client.Database(cfg.DatabaseName).Collection(cfg.UserCollection)
	filter := bson.M{"username": claims.Username}
	result := models.User{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.String(http.StatusBadRequest, "models.User not found")
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
	return
}

// TODO: remove one check login function here, and in handlers.go

// Check login or not
// Return a claim with username
func checkLogin_(c *gin.Context, client *mongo.Client) error {
	userID := c.Param("userID")
	// buf := make([]byte, 1024)
	// num, _ := c.Request.Body.Read(buf)
	// reqBody := string(buf[0:num])
	// fmt.Println(reqBody)
	// fmt.Println("id received", userID)
	return nil
	if userID == "" {
		return errors.New("You need to create an account before doing that!")
	}
	filter := bson.M{"_id": userID}
	status, msg, _ := checkForUser(filter, client, userID)
	fmt.Println(status)
	if status != http.StatusOK {
		return errors.New(msg)
	}
	return nil
}

func checkLogin(c *gin.Context) (*models.Claims, error) {
	t, err := c.Cookie("token")
	// Initialize a new instance of `models.Claims`
	claims := &models.Claims{}

	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			c.String(http.StatusUnauthorized, err.Error())
			return claims, err
		}
		// For any other type of error, return a bad request status
		c.String(http.StatusBadRequest, err.Error())
		return claims, err
	}

	// Get the JWT string from the cookie
	tknStr := t

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.String(http.StatusUnauthorized, err.Error())
			return claims, err
		}
		c.String(http.StatusBadRequest, err.Error())
		return claims, err
	}
	if !tkn.Valid {
		c.String(http.StatusUnauthorized, "Token is not valid!")
		return claims, err
	}
	return claims, err
}

// checkName checks the login stage and
// if the username outer handler gets is the same as the one in cookie
// usage:
// ok, err := checkName(c, username)
// if err != nil {
//// handle internal error of jwt
// c.String(http.StatusInternalServerError, err.Error())
// return
//}
// if !ok {
//// username not match
// c.String(http.StatusBadRequest, "models.User name not match!")
//}
func checkName(c *gin.Context, username string) (bool, error) {
	claims, err := checkLogin(c)
	if err != nil {
		return false, err
	} else if claims.Username != username {
		return false, nil
	}
	return true, nil
}
