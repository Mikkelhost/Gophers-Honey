package database

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

var db *mongo.Client

var (
	DB_NAME = "honeyDB"
	DB_HOST = getenv("DB_HOST", "cluster0.sb5ex.mongodb.net")
	DB_USER = getenv("DB_USER", "goadmin")
	DB_PASS = getenv("DB_PASS", "vcSXbkA7pBNbKpE8")
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().Msgf("Env %s not set, using default of %s", key, fallback)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func testInsertDocument() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coll := db.Database(DB_NAME).Collection("device_collection")
	_, err := coll.InsertOne(ctx, Device{
		UUID:  0x000000FF,
		IP:    0xFFFFFF00,
		IpStr: "255.255.255.0",
		Services: Service{
			SSH: true,
			FTP: true,
		},
	})
	if err != nil {
		return
	}
}

func Connect() {
	URI := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s", DB_USER, DB_PASS, DB_HOST, DB_NAME)
	clientOptions := options.Client().
		ApplyURI(URI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal().Msgf("Error connecting to DB: %s", err)
	}
	db = client
	testInsertDocument()
}

func Disconnect() {
	db.Disconnect(context.Background())
}
