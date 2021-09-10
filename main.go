package main

import (
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	"github.com/Mikkelhost/Gophers-Honey/pkg/httpserver"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

var DEBUG = true

func hashAndSaltPassword(pwd []byte) string {
	cost := 14
	hash, err := bcrypt.GenerateFromPassword(pwd, cost)
	if err != nil {
		log.Logger.Warn().Msgf("Error generating hash")
	}

	return string(hash)
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

	if database.IsUserInCollection(username, "username", database.DB_USER_COLL) {
		hashedPwd, err := database.GetPasswordHash(username)
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

func main() {
	// Initialize logger and set logging level.
	log.InitLog(DEBUG)

	// Set up database connection.
	log.Logger.Info().Msg("Setting up database connection")
	database.Connect()
	defer database.Disconnect()

	// Set up server.
	log.Logger.Info().Msg("Running server")
	httpserver.RunServer()
}
