package contracts

// AccountCreatedEvent is emitted whenever a new account is created
type AccountCreatedEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// EventName returns the event's name
func (c *AccountCreatedEvent) EventName() string {
	return "accountCreated"
}
