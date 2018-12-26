package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fintech-asean/fintech-go/account/msgqueue"
	"github.com/fintech-asean/fintech-go/persistence"
)

type accountHandler struct {
	dbHandler    persistence.DatabaseHandler
	eventEmitter msgqueue.EventEmitter
}

// newEventHandler initialises an accounthandler
func newAccountHandler(databaseHandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) *accountHandler {
	return &accountHandler{
		dbHandler:    databaseHandler,
		eventEmitter: eventEmitter,
	}
}

// allAccountHandler gets all the accounts from the DB
func (ah *accountHandler) allAccountHandler(w http.ResponseWriter, r *http.Request) {
	accounts, err := ah.dbHandler.FindAllAccounts()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error occured while trying to find all available accounts %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err = json.NewEncoder(w).Encode(&accounts)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error occured while trying encode account to JSON %s", err)
		return
	}

	w.WriteHeader(200)
}

// newAccountHandler saves the account to the DB
func (ah *accountHandler) newAccountHandler(w http.ResponseWriter, r *http.Request) {
	account := persistence.Account{}
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, fmt.Sprintf("Error occured while decoding account data: %v", err))
		return
	}

	_, err = ah.dbHandler.AddAccount(account)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error occured while persisting account %s", err)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&account)
}
