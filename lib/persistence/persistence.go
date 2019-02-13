package persistence

// DatabaseHandler is the handler that deals with accounts database
type DatabaseHandler interface {
	FindAllAccounts() ([]Account, error)
	AddAccount(Account) ([]byte, error)
	FindAccountByID([]byte) (Account, error)
}
