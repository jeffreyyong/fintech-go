package kafka

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/fintechasean/fintech-go/lib/msgqueue"
	log "github.com/sirupsen/logrus"
	"gopkg.in/Shopify/sarama.v1"
)

type kafkaEventListener struct {
	consumer   sarama.Consumer
	partitions []int32
	mapper     msgqueue.EventMapper
}

// NewKafkaListenerFromEnvironment creates a new Kafka Listener from the Environment
func NewKafkaListenerFromEnvironment() (msgqueue.EventListener, error) {
	brokers := []string{"localhost:9092"}
	partitions := []int32{}

	if brokerList := os.Getenv("KAFKA_BROKERS"); brokerList != "" {
		brokers = strings.Split(brokerList, ",")
	}

	if partitionList := os.Getenv("KAFKA_PARTITIONS"); partitionList != "" {
		partitionStrings := strings.Split(partitionList, ",")
		partitions = make([]int32, len(partitionStrings))

		for i := range partitionStrings {
			partition, err := strconv.Atoi(partitionStrings[i])
			if err != nil {
				return nil, err
			}
			partitions[i] = int32(partition)
		}
	}

	client := <-RetryReconnect(brokers, 5*time.Second)

	return NewKafkaEventListener(client, partitions)
}

// NewKafkaEventListener creates a new Kafka event listener
func NewKafkaEventListener(client sarama.Client, partitions []int32) (msgqueue.EventListener, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	listener := &kafkaEventListener{
		consumer:   consumer,
		partitions: partitions,
		mapper:     msgqueue.NewEventMapper(),
	}

	return listener, nil
}

func (k *kafkaEventListener) Listen(events ...string) (<-chan msgqueue.Event, <-chan error, error) {
	var err error

	topic := "accounts"
	results := make(chan msgqueue.Event)
	errors := make(chan error)

	partitions := k.partitions
	if len(partitions) == 0 {
		partitions, err = k.consumer.Partitions(topic)
		if err != nil {
			return nil, nil, err
		}
	}
	log.Infof("topic %s has partitions: %v", topic, partitions)

	for _, partition := range partitions {
		log.Infof("consuming partitions %s:%d", topic, partition)

		pConsumer, err := k.consumer.ConsumePartition(topic, partition, 0)
		if err != nil {
			return nil, nil, err
		}

		go func() {
			for msg := range pConsumer.Messages() {
				log.Infof("received message %v", msg)

				body := messageEnvelope{}
				err := json.Unmarshal(msg.Value, &body)
				if err != nil {
					errors <- fmt.Errorf("could not JSON-decode message: %v", err)
					continue
				}

				event, err := k.mapper.MapEvent(body.EventName, body.Payload)
				if err != nil {
					errors <- fmt.Errorf("could not map message: %v", err)
				}
				results <- event
			}
		}()

		go func() {
			for err := range pConsumer.Errors() {
				errors <- err
			}
		}()
	}

	return results, errors, nil
}
