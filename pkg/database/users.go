package database

import (
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// TODO: ADD error handling for usernames (empty string "", special chars, etc.).

type User struct {
	Username      string `bson:"username"`
	UsernameLower string `bson:"username_lower"`
	PasswordHash  string `bson:"password_hash"`
}

// AddNewUser adds a new user, with a specified username, to the database.
// TODO: HANDLE password info, salt and hash info when adding user.
func AddNewUser(username, hashedAndSaltedPwd string) {
	if IsUserInCollection(username, "username_lower", DB_USER_COLL) {
		log.Logger.Warn().Str("username", username).Msgf("Username already in use")
		return
	}

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	user := User{
		Username:      username,
		UsernameLower: strings.ToLower(username),
		PasswordHash:  hashedAndSaltedPwd,
	}

	_, err := db.Database(DB_NAME).Collection(DB_USER_COLL).InsertOne(ctx, user)

	if err != nil {
		log.Logger.Warn().Msgf("Error adding username: %s", err)
		return
	}
}

// IsUserInCollection reports whether a document with the specified
// username occurs in the given collection.
// TODO: CHECK if this method can be combined with isDeviceInCollection as they both do the same.
func IsUserInCollection(value, key, collection string) bool {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	filter := bson.M{
		key: strings.ToLower(value),
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

// RemoveUser removes a user, with the specified username, from the
// database.
func RemoveUser(username string) {
	if !IsUserInCollection(username, "username_lower", DB_USER_COLL) {
		log.Logger.Warn().Str("username", username).Msgf("Username not found")
		return
	}

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	filter := bson.M{
		"username": username,
	}

	_, err := db.Database(DB_NAME).Collection(DB_USER_COLL).DeleteOne(ctx, filter)

	if err != nil {
		log.Logger.Warn().Msgf("Error removing user: %s", err)
		return
	}
}

// GetPasswordHash retrieves the stored password hash for the specified
// username.
func GetPasswordHash(username string) (string, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()
	log.Logger.Debug().Msgf("Getting hash for username: %s", username)
	filter := User{
		UsernameLower: strings.ToLower(username),
	}

	var user User

	result := db.Database(DB_NAME).Collection(DB_USER_COLL).FindOne(ctx, filter)

	if err := result.Decode(&user); err != nil {
		//log.Logger.Warn().Msgf("Error decoding result: %s", err)
		return "", err
	}

	return user.PasswordHash, nil
}

// verifyPassword compares a plaintext password with a hashed and salted
// password and returns true if they match
func verifyPassword(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)

	if err != nil {
		log.Logger.Warn().Str("hash", hashedPwd).Msgf("Password does not match hash")
		return false
	}
	return true
}

// loginUser verifies a user login by checking whether the password
func loginUser(username, stringPwd string) (bool, error) {
	pwd := []byte(stringPwd)

	if IsUserInCollection(username, "username", DB_USER_COLL) {
		hashedPwd, err := GetPasswordHash(username)
		if err != nil {
			log.Logger.Warn().Msgf("Error retrieving password hash")
			return false, err
		}
		verified := verifyPassword(hashedPwd, pwd)
		if !verified {
			return false, nil
		}
		return true, nil
	}
	return false, nil
}
