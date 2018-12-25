package msgqueue

// Event is the interface definition for events that are emitted using an EventEmitter
type Event interface {
	Name() string
}
