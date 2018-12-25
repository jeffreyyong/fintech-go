package kafka

import (
	"encoding/json"

	"github.com/fintech-asean/fintech-go/account/msgqueue"
	sarama "gopkg.in/Shopify/sarama.v1"
)

type kafkaEventEmitter struct {
	producer sarama.SyncProducer
}

type messageEnvelope struct {
	EventName string      `json:"eventName"`
	Payload   interface{} `json:"payload"`
}

// NewKafkaEventEmitter initialises a new Kafka Event Emitter
func NewKafkaEventEmitter(client sarama.Client) (msgqueue.EventEmitter, error) {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	emitter := kafkaEventEmitter{
		producer: producer,
	}

	return &emitter, nil
}

// Emit emit events
func (k *kafkaEventEmitter) Emit(evt msgqueue.Event) error {
	jsonBody, err := json.Marshal(messageEnvelope{
		evt.Name(),
		evt,
	})

	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "accounts",
		Value: sarama.ByteEncoder(jsonBody),
	}

	_, _, err = k.producer.SendMessage(msg)

	return err
}
