package database

import (
	"errors"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

// TODO: ADD error handling for usernames (empty string "", special chars, etc.).

// AddNewUser adds a new user, with a specified username, to the database.
func AddNewUser(user model.DBUser, hashedAndSaltedPwd string) error {
	if IsUserInCollection(strings.ToLower(user.Username), "username_lower", DB_USER_COLL) {
		log.Logger.Warn().Str("username", user.Username).Msgf("Username already in use")
		return errors.New("username already exists")
	}

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	dbUser := model.DBUser{
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Username:      user.Username,
		UsernameLower: strings.ToLower(user.Username),
		PasswordHash:  hashedAndSaltedPwd,
		Role:          user.Role,
	}

	_, err := db.Database(DB_NAME).Collection(DB_USER_COLL).InsertOne(ctx, dbUser)

	if err != nil {
		log.Logger.Warn().Msgf("Error adding username: %s", err)
		return err
	}
	return nil
}

// IsUserInCollection reports whether a document with the specified
// username occurs in the given collection.
// TODO: CHECK if this method can be combined with isIdInCollection as they both do the same.
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
func RemoveUser(username string) error{
	if !IsUserInCollection(strings.ToLower(username), "username_lower", DB_USER_COLL) {
		log.Logger.Warn().Str("username", username).Msgf("Username not found")
		return errors.New("Username not in db")
	}

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	filter := bson.M{
		"username_lower": strings.ToLower(username),
	}

	_, err := db.Database(DB_NAME).Collection(DB_USER_COLL).DeleteOne(ctx, filter)

	if err != nil {
		log.Logger.Warn().Msgf("Error removing user: %s", err)
		return err
	}
	return nil
}

// UpdateUser updates user-specific fields in the user collection.
func UpdateUser(user model.APIUser, hashedAndSaltedPwd string) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	filter := bson.M{
		"username_lower": strings.ToLower(user.Username),
	}

	conf := bson.M{
		"notifications_enabled": user.NotificationsEnabled,
	}

	// If user set a new email, update field.
	if user.Email != "" {
		conf["email"] = user.Email
	}

	// If user set a new password, update field.
	if hashedAndSaltedPwd != "" {
		conf["password_hash"] = hashedAndSaltedPwd
	}

	update := bson.M{
		"$set": conf,
	}

	_, err := db.Database(DB_NAME).Collection(DB_USER_COLL).UpdateOne(ctx, filter, update)
	if err != nil {
		log.Logger.Warn().Msgf("Error updating user: %s", err)
		return err
	}

	return nil
}

// GetPasswordHash retrieves the stored password hash for the specified
// username.
func GetPasswordHash(username string) (string, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()
	log.Logger.Debug().Msgf("Getting hash for username: %s", username)

	filter := bson.M{
		"username_lower": strings.ToLower(username),
	}

	var user model.DBUser

	result := db.Database(DB_NAME).Collection(DB_USER_COLL).FindOne(ctx, filter)

	if err := result.Decode(&user); err != nil {
		log.Logger.Warn().Msgf("Error decoding result: %s", err)
		return "", err
	}

	return user.PasswordHash, nil
}

// LoginUser verifies a user login by checking whether provided password
// matches the hashed password stored under the specified username.
func LoginUser(username, stringPwd string) (bool, error) {
	pwd := []byte(stringPwd)

	if IsUserInCollection(username, "username_lower", DB_USER_COLL) {
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
	} else {
		// If the username is not valid we do a faux password hash compare
		// in order for attackers not to be able to enumerate usernames by
		// timing hash compare time.
		// TODO: Important to mention in report
		dummyHash := "$2a$14$V4MAXIGnk26YP9xOlhxUn.PW45vqUzLtoE4eGz0TD1m1R6i6IcMEq"
		_ = verifyPassword(dummyHash, pwd)
		return false, nil
	}
}

func GetUser(username string) (model.DBUser, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	filter := bson.M{
		"username_lower": strings.ToLower(username),
	}
	result := db.Database(DB_NAME).Collection(DB_USER_COLL).FindOne(ctx, filter)

	var user model.DBUser

	err := result.Decode(&user)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding result: %s", err)
		return model.DBUser{}, err
	}
	user.PasswordHash = ""
	return user, nil
}

// GetAllUsers retrieves all users currently in the database,
// and removes the hashed password of the users before returning the information
func GetAllUsers() ([]model.DBUser, error) {
	var userList []model.DBUser

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	results, err := db.Database(DB_NAME).Collection(DB_USER_COLL).Find(ctx, bson.M{})

	if err != nil {
		log.Logger.Warn().Msgf("Error retrieving user list: %s", err)
		return nil, err
	}

	for results.Next(ctx) {
		var user model.DBUser

		if err = results.Decode(&user); err != nil {
			log.Logger.Warn().Msgf("Error decoding result: %s", err)
			return nil, err
		}
		user.PasswordHash = ""
		userList = append(userList, user)
	}

	for _, user := range userList {
		log.Logger.Debug().Msgf("Found user with user ID: %d", user.UsernameLower)
	}

	return userList, nil
}
