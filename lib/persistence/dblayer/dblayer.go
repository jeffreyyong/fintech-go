package dblayer

import (
	"github.com/jeffreyyong/fintech-go/lib/persistence"
	"github.com/jeffreyyong/fintech-go/lib/persistence/mongo"
)

// DbType is database type
type DbType string

const (
	// MongoDB enum
	MongoDB DbType = "mongodb"
)

// NewPersistenceLayer initialises a new persistence layer
func NewPersistenceLayer(options DbType, connection string) (persistence.DatabaseHandler, error) {
	switch options {
	case MongoDB:
		return mongo.NewMongoDB(connection)
	}
	return nil, nil
}
