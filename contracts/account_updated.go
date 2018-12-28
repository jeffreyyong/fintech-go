package contracts

// AccountUpdatedEvent is emitted whenever an account is updated
type AccountUpdatedEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// EventName returns the event's name
func (c *AccountUpdatedEvent) EventName() string {
	return "accountUpdated"
}
