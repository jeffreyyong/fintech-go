package msgqueue

import (
	"encoding/json"
	"fmt"

	"github.com/jeffreyyong/fintech-go/contracts"
	"github.com/mitchellh/mapstructure"
)

// EventMapper is an interface that maps string to Event
type EventMapper interface {
	MapEvent(string, interface{}) (Event, error)
}

// NewEventMapper initialises a new EventMapper
func NewEventMapper() EventMapper {
	return &eventMapper{}
}

type eventMapper struct{}

func (eventMapper *eventMapper) MapEvent(eventName string, serialized interface{}) (Event, error) {
	var event Event

	switch eventName {
	case "accountUpdated":
		event = &contracts.AccountUpdatedEvent{}
	case "accountCreated":
		event = &contracts.AccountCreatedEvent{}
	default:
		return nil, fmt.Errorf("unknown event type %s", eventName)
	}

	switch s := serialized.(type) {
	case []byte:
		err := json.Unmarshal(s, event)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal event %s: %s", eventName, err)
		}
	default:
		cfg := mapstructure.DecoderConfig{
			Result:  event,
			TagName: "json",
		}
		dec, err := mapstructure.NewDecoder(&cfg)
		if err != nil {
			return nil, fmt.Errorf("could not initialize decoder for event %s: %s", eventName, err)
		}
		err = dec.Decode(s)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal event %s: %s", eventName, err)
		}
	}

	return event, nil
}
