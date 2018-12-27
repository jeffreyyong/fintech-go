package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/fintechasean/fintech-go/account/msgqueue"
	"bitbucket.org/fintechasean/fintech-go/persistence"
	"github.com/go-chi/chi"
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
		fmt.Fprintf(w, "Error occured while trying encode accounts to JSON %s", err)
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

// oneAccountHandler gets a specific account from the DB based on account ID
func (ah *accountHandler) oneAccountHandler(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")

	if accountID == "" {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Missing route parameter accountID")
		return
	}

	accountIDBytes, err := hex.DecodeString(accountID)

	account, err := ah.dbHandler.FindAccountByID(accountIDBytes)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "account with id %s was not found, err: %v", accountID, err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err = json.NewEncoder(w).Encode(&account)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error occured while trying to encode account to JSON %s", err)
		return
	}

	w.WriteHeader(200)
}
