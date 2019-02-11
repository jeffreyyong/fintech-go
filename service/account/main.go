package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/fintechasean/fintech-go/configuration"
	"bitbucket.org/fintechasean/fintech-go/lib/msgqueue"
	"bitbucket.org/fintechasean/fintech-go/lib/persistence/dblayer"
	"bitbucket.org/fintechasean/fintech-go/service/account/rest"
	"github.com/go-kit/kit/log"
	// "github.com/go-kit/kit/log"
)

func main() {
	// create logger
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// load configuration
	confPath := flag.String("conf", "../../configuration/config.yaml", "flag to set the path to the configuration JSON file")
	flag.Parse()
	// extract configuration
	config, err := configuration.ExtractConfiguration(*confPath)
	if err != nil {
		logger.Log("Error", "Fail to extract configuration")
		panic(err)
	}

	// create event emitter
	var eventEmitter msgqueue.EventEmitter
	// conf := sarama.NewConfig()
	// conf.Producer.Return.Successes = true
	// conn, err := sarama.NewClient(config.KafkaMessageBrokers, conf)
	// if err != nil {
	// 	panic(err)
	// }

	// eventEmitter, err = kafka.NewKafkaEventEmitter(conn)
	// if err != nil {
	// 	panic(err)
	// }

	// create dbHandler
	dbHandler, err := dblayer.NewPersistenceLayer(config.DBType, config.DBConnection)
	if err != nil {
		logger.Log("Error", "Fail to initialise a dbhandler")
		panic(err)
	}

	srv := rest.New(dbHandler, eventEmitter, log.With(logger, "component", "http"))

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", config.RestEndpoint, "msg", "listening")
		errs <- http.ListenAndServe(config.RestEndpoint, srv)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	logger.Log("terminated", <-errs)
}
