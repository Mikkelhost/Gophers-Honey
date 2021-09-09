package database

import (
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: ADD error handling for usernames (empty string "", special chars, etc.).
// TODO: ADD functionality for case-sensitive usernames and/or username checks.

type User struct {
	USERNAME string `bson:"username"`
}

// AddNewUser adds a new user, with a specified username to the database.
// TODO: HANDLE password info, salt and hash info when adding user.
func AddNewUser(username, salt, hash string) {
	if isUserInCollection(username, "username", DB_USER_COLL) {
		log.Logger.Fatal().Msgf("Username already in use")
		return
	}

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	user := User{
		USERNAME: username,
	}

	_, err := db.Database(DB_NAME).Collection(DB_USER_COLL).InsertOne(ctx, user)

	if err != nil {
		log.Logger.Fatal().Msgf("Error adding username: %s", err)
		return
	}
}

// TODO: CHECK if this method can be combined with isDeviceInCollection as they both do the same.
// isUserInCollection reports whether a document with the specified
// username occurs in the given collection.
func isUserInCollection(value, key, collection string) bool {
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
