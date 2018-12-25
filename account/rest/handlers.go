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
		fmt.Fprintf(w, "Error occured whilew trying encode events to JSON %s", err)
	}
}
