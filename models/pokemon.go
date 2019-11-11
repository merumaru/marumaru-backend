package models

import (
	"time"

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
	CreateData time.Time
}

type Product struct {
	ID          primitive.ObjectID `bson:"_id, omitempty"`
	Photos      []string
	SellerName  string
	Name        string
	Description string
	Price       float32
	TimeDuration
	Tags []string
}
