package mongo

import (
	"bitbucket.org/fintechasean/fintech-go/persistence"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	// DbName is the database name for myaccounts
	DbName    = "bankaccounts"
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

// FindAccountByID finds an account by ID
func (mg *MongoDB) FindAccountByID(ID []byte) (persistence.Account, error) {
	return persistence.Account{}, nil
}

// AddAccount adds account to the DB
func (mg *MongoDB) AddAccount(account persistence.Account) ([]byte, error) {
	s := mg.getFreshSession()
	defer s.Close()

	if !account.ID.Valid() {
		account.ID = bson.NewObjectId()
	}

	for i, transaction := range account.Transactions {
		if transaction.ID == "" {
			account.Transactions[i].ID = bson.NewObjectId()
		}
	}

	return []byte(account.ID), s.DB(DbName).C(Accounts).Insert(account)
}

func (mg *MongoDB) getFreshSession() *mgo.Session {
	return mg.session.Copy()
}
