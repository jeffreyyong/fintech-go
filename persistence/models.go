package persistence

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// Account is the bank account that the user has
type Account struct {
	ID           bson.ObjectId `bson:"_id"`
	Name         string
	Balance      float64
	AccountLogo  string
	LastUpdated  time.Duration
	Transactions []Transaction
}

// Transaction holds the transaction details
type Transaction struct {
	ID           bson.ObjectId `bson:"_id"`
	Name         string
	Date         time.Duration
	MerchantName string
	MerchantLogo string
	Amount       float64
	Categories   []Category
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
