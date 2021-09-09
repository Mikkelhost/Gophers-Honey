package database

import (
	"context"
	"fmt"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var db *mongo.Client

var (
	DB_NAME      = "honeyDB"
	DB_HOST      = getenv("DB_HOST", "cluster0.sb5ex.mongodb.net")
	DB_USER      = getenv("DB_USER", "goadmin")
	DB_PASS      = getenv("DB_PASS", "vcSXbkA7pBNbKpE8")
	DB_DEV_COLL  = "device_collection"
	DB_CONF_COLL = "config_collection"
	DB_USER_COLL = "user_collection"
	DEBUG        = false
)

// getenv retrieves the value of the environment variable named by the
// key. If no environment variable of the provided key is found a
// fallback is used as a default value.
func getenv(key, fallback string) string {
	value := os.Getenv(key)
	log.Logger.Debug().Msgf("Env %s not set, using default of %s", key, fallback)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// Connect creates a connection to the database.
func Connect() {
	URI := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s", DB_USER, DB_PASS, DB_HOST, DB_NAME)
	clientOptions := options.Client().ApplyURI(URI)
	ctx, cancel := getContextWithTimeout()
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Logger.Fatal().Msgf("Error connecting to DB: %s", err)
	}
	db = client
}

// Disconnect shuts down the current database connection.
func Disconnect() {
	if db == nil {
		log.Logger.Warn().Msgf("No database connection to disconnect.")
		return
	}
	err := db.Disconnect(context.Background())
	if err != nil {
		log.Logger.Fatal().Msgf("Error disconnecting from DB: %s", err)
	}
}
