package server

import (
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/merumaru/marumaru-backend/data"
	models "github.com/merumaru/marumaru-backend/models"

	"go.mongodb.org/mongo-driver/mongo"
)

//signUp registers a new user
func signUp(c *gin.Context, client *mongo.Client) {
	var newuser models.User
	//TODO check for error
	if err := c.ShouldBindJSON(&newuser); err != nil {
		log.Println("Error in binding JSON when signing up : ", err.Error())
		c.String(500, fmt.Sprintf("Error occured while signing up"))
		return
	}

	if err := data.AddUser(client, &newuser); err != nil {
		c.String(400, err.Error())
	} else {
		c.String(200, fmt.Sprintf("User Successfully added"))
	}
}

func isUserSignedIn(c *gin.Context) bool {
	if u := getUserByCookie(c); u != nil {
		return true
	}
	return false
}

// GetUserByCookie returns the whole user struct by cookie
func getUserByCookie(c *gin.Context) *models.User {
	if jwtToken, err := authMiddleware.ParseToken(c); err == nil {
		if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
			if userID, ok := claims[jwtUserNameKey].(string); ok == true {
				if user, _ := data.GetUserByUserName(mongoDBClient, userID); user != nil {
					log.Println("User is signed in")
					return user
				}
			} else {
				log.Printf("Could not convert claim to string")
			}
		} else {
			log.Printf("Could not extract claims into JWT")
		}
	}
	log.Println("User is not signed in")
	return nil
}
