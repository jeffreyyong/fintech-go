package persistence

import (
	"github.com/globalsign/mgo/bson"
)

// Account is the bank account that the user has
type Account struct {
	ID           bson.ObjectId `bson:"_id"`
	Name         string        `json:"name"`
	Balance      float64       `json:"balance"`
	AccountLogo  string        `json:"account_logo"`
	LastUpdated  string        `json:"last_updated"`
	Transactions []Transaction `json:"transactions"`
}

// Transaction holds the transaction details
type Transaction struct {
	ID           bson.ObjectId `bson:"_id"`
	Name         string        `json:"name"`
	Date         string        `json:"date"`
	MerchantName string        `json:"merchant_name"`
	MerchantLogo string        `json:"merchant_logo"`
	Amount       float64       `json:"amount"`
	Category     Category      `json:"category"`
}

// Category is the type for category
type Category int

const (
	Transfers Category = iota
	General
	Housing
	Shopping
	EatingOut
	Travel
	Groceries
	Leisure
	Income
)

func (c Category) String() string {
	return [...]string{"Transfers", "General", "Housing", "Shopping", "Eating out", "Travel", "Groceries", "Leisure", "Income"}[c]
}
