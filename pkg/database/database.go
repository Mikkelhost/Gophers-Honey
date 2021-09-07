package database

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func getenv(key, fallback string) string {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	value := os.Getenv(key)
	log.Debug().Msgf("Env %s not set, using default of %s", key, fallback)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func Connect() {
	URI := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s", DB_USER, DB_PASS, DB_HOST, DB_NAME)
	clientOptions := options.Client().
		ApplyURI(URI)
	ctx, cancel := getContextWithTimeout()
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal().Msgf("Error connecting to DB: %s", err)
	}
	db = client
	//ConfigureDevice(Service{RDP: true, FTP: true}, 3311712553)
	//AddDevice("10.0.0.3")
	_ = GetAllDevices()
}

func Disconnect() {
	db.Disconnect(context.Background())
}
