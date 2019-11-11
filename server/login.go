package server

import (
	"fmt"
	"net/http"
	"time"

	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

// Create handler which needs db client
type HandlerFuncWithDB func(*gin.Context, *mongo.Client)

// Suppose users just here, add database later
// var users = map[string]string{
// 	"user1": "password1",
// 	"user2": "password2",
// }

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Signin add a user cookie
func Signin(c *gin.Context, client *mongo.Client) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.String(400, err.Error())
		return
	}

	collection := client.Database("testing").Collection("users")
	filter := bson.M{"username": creds.Username}
	result := User{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		// No result
		if err == mongo.ErrNoDocuments {
			c.String(http.StatusUnauthorized, fmt.Sprintf("No such user name %s!", creds.Username))
			return
		} else {
			// Other problem
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}
	// expectedPassword, ok := users[creds.Username]
	// Check password
	if result.Password != creds.Password {
		c.String(http.StatusUnauthorized, "Password is incorrected!")
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
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
	c.SetCookie("token", tokenString, expirationTime.Second(), "/", "", true, false)
	c.String(http.StatusOK, "Set cookie successfully!")
	return
}

// Signup adds a new user
func SignUp(c *gin.Context, client *mongo.Client) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.String(400, err.Error())
		return
	}
	// Check input
	if user.Username == "" || user.Password == "" || user.Email == "" {
		c.String(400, "Invalid input")
		return
	}
	collection := client.Database("testing").Collection("users")
	filter := bson.M{"username": user.Username}
	result := User{}
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
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "Sign up successfully!")
}

// Welcome tests login
func Welcome(c *gin.Context) {
	claims, err := checkLogin(c)
	if err != nil {
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("Welcome %s!", claims.Username))
}

// Refresh the token in the background by the client application.
func Refresh(c *gin.Context) {
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

	c.SetCookie("token", tokenString, expirationTime.Second(), "/", "", true, false)
	c.String(http.StatusOK, "Refresh a new token!")
	return
}

// Check login or not
func checkLogin(c *gin.Context) (*Claims, error) {
	t, err := c.Cookie("token")
	// Initialize a new instance of `Claims`
	claims := &Claims{}

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
// c.String(http.StatusBadRequest, "User name not match!")
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
