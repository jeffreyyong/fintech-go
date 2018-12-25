package mongo

import (
	"github.com/fintech-asean/fintech-go/persistence"
	"github.com/globalsign/mgo"
)

const (
	// DbName is the database name for myaccounts
	DbName    = "myaccounts"
	Customers = "customers"
	Accounts  = "accounts"
	Balances  = "balances"
)

// MongoDB is a struct that holds the mongoDB session
type MongoDB struct {
	session *mgo.Session
}

// NewMongoDB initialises a new mongoDB instance
func NewMongoDB(connection string) (persistence.DatabaseHandler, error) {
	s, err := mgo.Dial(connection)
	return &MongoDB{
		session: s,
	}, err
}

// FindAllAccounts lists all the accounts
func (mg *MongoDB) FindAllAccounts() ([]persistence.Account, error) {
	s := mg.getFreshSession()
	defer s.Close()
	accounts := []persistence.Account{}
	err := s.DB(DbName).C(Accounts).Find(nil).All(&accounts)
	return accounts, err
}

func (mg *MongoDB) getFreshSession() *mgo.Session {
	return mg.session.Copy()
}
