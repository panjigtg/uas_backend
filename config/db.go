package config

import (
	"uas/database"

	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseContainer struct {
	Postgres *sql.DB
	Mongo    *mongo.Database
}

func InitDatabase() *DatabaseContainer {
	return &DatabaseContainer{
		Postgres: database.PostgresConnections(),
		Mongo:    database.MongoConnections(),
	}
}
