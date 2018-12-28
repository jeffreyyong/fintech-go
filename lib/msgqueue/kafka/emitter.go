package kafka

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"bitbucket.org/fintechasean/fintech-go/lib/msgqueue"
	log "github.com/sirupsen/logrus"
	sarama "gopkg.in/Shopify/sarama.v1"
)

type kafkaEventEmitter struct {
	producer sarama.SyncProducer
}

type messageEnvelope struct {
	EventName string      `json:"eventName"`
	Payload   interface{} `json:"payload"`
}

// NewKafkaEventEmitterFromEnvironment returns a new Kafka Event Emitter from the environment
func NewKafkaEventEmitterFromEnvironment() (msgqueue.EventEmitter, error) {
	brokers := []string{"localhost:9092"}

	if brokerList := os.Getenv("KAFKA_BROKERS"); brokerList != "" {
		brokers = strings.Split(brokerList, ",")
	}

	client := <-RetryReconnect(brokers, 5*time.Second)
	return NewKafkaEventEmitter(client)
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
		evt.EventName(),
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

// RetryReconnect implements a retry mechanism for establishing the Kafka connection.
// This is necessary in container environments where individual components may be started out-of-order.
// Might have to wait for upstream services like Kafka to actually become available.
//
// Alternatives:
// 	- use an entrypoint script in the container in which you wait for the service to be available
func RetryReconnect(brokers []string, retryInterval time.Duration) chan sarama.Client {
	result := make(chan sarama.Client)

	go func() {
		defer close(result)
		for {
			config := sarama.NewConfig()
			conn, err := sarama.NewClient(brokers, config)
			if err == nil {
				log.Info("connection successfully established")
				result <- conn
				return
			}

			log.Warnf("Kafka connection failed with error (retrying in %s): %s", retryInterval.String(), err)
			time.Sleep(retryInterval)
		}
	}()

	return result
}
