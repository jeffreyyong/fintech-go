package persistence

// DatabaseHandler is the handler that deals with accounts database
type DatabaseHandler interface {
	FindAllAccounts() ([]Account, error)
}
