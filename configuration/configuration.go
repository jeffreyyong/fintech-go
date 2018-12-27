package configuration

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"bitbucket.org/fintechasean/fintech-go/persistence/dblayer"
	yaml "gopkg.in/yaml.v2"
)

var (
	// DBTypeDefault is the default dblayer type
	DBTypeDefault             = dblayer.DbType("mongodb")
	DBConnectionDefault       = "localhost"
	RestEndpointDefault       = "localhost:8181"
	MessageBrokerTypeDefault  = "kafka"
	KafkaMessageBrokerDefault = []string{"localhost:9092"}
)

// ServiceConfig holds the config for the service
type ServiceConfig struct {
	DBType              dblayer.DbType `yaml:"db_type"`
	DBConnection        string         `yaml:"db_connection"`
	RestEndpoint        string         `yaml:"rest_endpoint"`
	MessageBrokerType   string         `yaml:"message_broker_type"`
	KafkaMessageBrokers []string       `yaml:"kafka_messsage_brokers"`
}

// ExtractConfiguration takes the file name and returns the ServiceConfig
func ExtractConfiguration(fileName string) (*ServiceConfig, error) {
	cfg := &ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestEndpointDefault,
		MessageBrokerTypeDefault,
		KafkaMessageBrokerDefault,
	}

	source, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(source, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.DBType == "" {
		return nil, fmt.Errorf("Missing db type")
	}

	if cfg.DBConnection == "" {
		return nil, fmt.Errorf("Missing db connection")
	}

	if cfg.RestEndpoint == "" {
		return nil, fmt.Errorf("Missing rest endpoint")
	}

	if cfg.MessageBrokerType == "" {
		return nil, fmt.Errorf("Missing broker type")
	}

	if reflect.DeepEqual(cfg.KafkaMessageBrokers, []string{}) {
		return nil, fmt.Errorf("Missing kafka message brokers")
	}

	return cfg, nil
}
