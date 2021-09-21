package database

import (
	"bytes"
	"context"
	"encoding/binary"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net"
	"time"
)

// createRandID pseudo-randomly generates a number which is checked
// against the device IDs currently in the collection. Returns when no
// collision is detected.
func createRandID(key, collection string) uint32 {
	rand.Seed(time.Now().Unix())
	deviceID := rand.Uint32()
	for isIdInCollection(deviceID, key, collection) {
		deviceID = rand.Uint32()
		log.Logger.Debug().Msg("Running \"while loop\"") //TODO: No need to keep this?
	}
	return deviceID
}

// isIdInCollection reports whether a document with the specified
// device ID occurs in the given collection.
func isIdInCollection(value uint32, key, collection string) bool {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	filter := bson.M{
		key: value,
	}

	countOptions := options.Count().SetLimit(1)
	count, err := db.Database(DB_NAME).Collection(collection).CountDocuments(ctx, filter, countOptions)

	if err != nil {
		log.Logger.Warn().Msgf("Error counting documents: %s", err)
	}

	if count > 0 {
		return true
	}

	return false
}

// ip2int converts an IP address from its string representation to its
// integer value.
func ip2int(ipStr string) uint32 {
	var long uint32
	err := binary.Read(bytes.NewBuffer(net.ParseIP(ipStr).To4()), binary.BigEndian, &long)
	if err != nil {
		log.Logger.Warn().Msgf("Error converting IP to int: %s", err)
		return 0
	}
	return long
}

// getContextWithTimeout is used to get a timeout context used when
// communicating with MongoDB.CompareHashAndPassword
func getContextWithTimeout() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx, cancel
}

// verifyPassword compares a plaintext password with a hashed and salted
// password and returns true if they match
func verifyPassword(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)

	if err != nil {
		log.Logger.Warn().Str("hash", hashedPwd).Msg("Password does not match hash")
		return false
	}
	return true
}

// stringAppearsInArray checks whether a given string occurs in an array.
func stringAppearsInArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// getIndexNames returns the index names of a given collection.
func getIndexNames(collection string) ([]string, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	var indexNames []string
	results, _ := db.Database(DB_NAME).Collection(collection).Indexes().List(ctx)

	for results.Next(ctx) {
		var indexName mongo.IndexSpecification

		err := results.Decode(&indexName)

		indexNames = append(indexNames, indexName.Name)
		if err != nil {
			log.Logger.Warn().Msgf("Error decoding indexes: %s", err)
			return nil, nil
		}
	}
	return indexNames, nil
}
