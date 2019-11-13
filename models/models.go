package models

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Book = iota
	Game = iota
	Toy  = iota
)

type TimeDuration struct {
	Start time.Time
	End   time.Time
}

type Order struct {
	ID         primitive.ObjectID `bson:"_id, omitempty"`
	SellerName string
	BuyerName  string
	ProductID  primitive.ObjectID
	TimeDuration
	IsCancelled bool
}

type Product struct {
	ID          primitive.ObjectID `bson:"_id, omitempty"`
	Photos      []string           `json:"photos"`
	SellerID    string				`json:"userID"`
	Name        string				`json:"name"`
	Description string				`json:"description"`
	Price       float32				`json:"price"`
	TimeDuration					`json:"timeduration"`
	Tags []string					`json:"tags"`
}


// Suppose users just here, add database later
// var users = map[string]string{
// 	"user1": "password1",
// 	"user2": "password2",
// }

type User struct {
	ID          primitive.ObjectID `bson:"_id, omitempty"`
	Username    string             `json:"username"`
	Password    string             `json:"password"`
	Email       string             `json:"email"`
	Address     string             `json:"address"`
	PhoneNumber string             `json:"phonenumber"`
	Avatar 		string 				`json:"avatar"`
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