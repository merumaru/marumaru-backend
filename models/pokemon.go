package models

import "time"

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
	SellerID  string
	BuyerID   string
	ProductID string
	TimeDuration
	CreateData int
}

type Product struct {
	Photos      []string
	SellerID    string
	Name        string
	Description string
	Price       float32
	TimeDuration
	Tags []int
}
