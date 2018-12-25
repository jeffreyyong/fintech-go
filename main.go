package main

import (
	"flag"
	"os"

	"github.com/go-kit/kit/log"
)

// ability to build and execute tests in a very afst manatner
// balances
// transactions
// schema per service

const (
	defaultPort                  = "8080"
	defaultTransactionServiceURL = "http://localhost:8080"
	defaultMongoDBURL            = "127.0.0.1"
	defaultDBName                = "default_db_name"
)

func main() {
	var (
	// addr   = envString("PORT", defaultPort)
	// tsurl  = envString("TRANSACTION_SERVICE_URL", defaultTransactionServiceURL)
	// dbURL  = envString("MONGODB_URL", defaultMongoDBURL)
	// dbName = envString("DB_NAME", defaultDBName)

	// httpAddr              = flag.String("http.addr", ":"+addr, "HTTP listen address")
	// transactionServiceURL = flag.String("service.transaction", tsurl, "transactoin service URL")
	// mongoDBURL            = flag.String("db.url", dbURL, "MongoDB URL")
	// databaseName          = flag.String("db.name", dbName, "MongoDB database name")
	// inMemory              = flag.Bool("inmem", false, "use in-memory repositories")

	// ctx = context.Background()
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// var ts transaction.Service
	// tsr = transaction.NewService()

}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e

}
