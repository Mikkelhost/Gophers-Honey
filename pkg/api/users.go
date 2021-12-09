package api

import (
	"encoding/json"
	"fmt"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"net/http"
	"regexp"
	"strings"

	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/gorilla/mux"
)

/*
The users API handles everything about users

All functions should write json data to the responsewriter
*/

func usersSubrouter(r *mux.Router) {
	usersAPI := r.PathPrefix("/api/users").Subrouter()
	usersAPI.Queries("user", "{user:.+}").HandlerFunc(tokenAuthMiddleware(getUser)).Methods("GET", "OPTIONS").Name("user")
	usersAPI.HandleFunc("", tokenAuthMiddleware(userHandler)).Methods("GET", "POST", "PUT", "DELETE", "OPTIONS")
	usersAPI.HandleFunc("/login", loginUser).Methods("POST", "OPTIONS")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// CORS preflight handling.
	if r.Method == "OPTIONS" {
		return
	}
	switch r.Method {
	case "GET":
		getUsers(w, r)
		return
	case "POST":
		registerUser(w, r)
		return
	case "PUT":
		updateUser(w, r)
		return
	case "DELETE":
		deleteUser(w, r)
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["user"]
	regex := "^[a-zA-Z0-9]*$"
	log.Logger.Debug().Msgf("Checking if username matches regex: %s", strings.TrimSpace(username))
	found, err := regexp.Match(regex, []byte(strings.TrimSpace(username)))
	log.Logger.Debug().Bool("found", found).Msg("Found is")
	if !found {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Illegal characters in username"})
		return
	}

	log.Logger.Info().Str("username", username).Msg("Getting user")
	user, err := database.GetUser(username)
	if err != nil {
		log.Logger.Warn().Msgf("Error getting user from database: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Error getting user from db"})
		return
	}
	json.NewEncoder(w).Encode(user)
}

// getUsers gets all information of all users, but the hashed passwords
func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []model.DBUser
	users, err := database.GetAllUsers()
	log.Logger.Debug().Msgf("Users: %v", users)
	if err != nil {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Error retrieving users"})
		return
	}
	if len(users) == 0 {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "No users in DB"})
		return
	}
	json.NewEncoder(w).Encode(users)
}

// loginUser Authenticates a user by their username and password.
// It returns a JWT token for the user to use as a session cookie.
// The JWT will be used forward as authentication for authenticated API endpoints.
func loginUser(w http.ResponseWriter, r *http.Request) {
	var userInfo = model.APIUser{}
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInfo); err != nil {
		log.Logger.Warn().Msgf("Failed decoding json: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding json: %s", err)})
		return
	}
	log.Logger.Debug().Msgf("Username: %s, Password: %s", userInfo.Username, userInfo.Password)
	loginStatus, err := database.LoginUser(userInfo.Username, userInfo.Password)
	log.Logger.Debug().Bool("loginStatus", loginStatus).Msg("GetPasswordHash result")
	if err != nil {
		log.Logger.Warn().Msgf("Error Loggin in user: /%s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("%s", err)})
		return
	}
	if !loginStatus {
		log.Logger.Debug().Msg("Incorrect username or password")
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Incorrect username or password"})
		return
	}
	user, err := database.GetUser(userInfo.Username)
	if err != nil {
		log.Logger.Warn().Msgf("Error getting user: %s from DB: %s", userInfo.Username, err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Internal server error"})
		return
	}
	token, err := createToken(user)
	if err != nil {
		log.Logger.Warn().Msgf("Error creating token: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("%s", err)})
		return
	}
	json.NewEncoder(w).Encode(model.APIUser{Token: token})
}

// registerUser creates a new user from the given information,
// and hashes and salts password.
func registerUser(w http.ResponseWriter, r *http.Request) {
	var newUser = model.APIUser{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
		log.Logger.Warn().Msgf("Error decoding json: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding json: %s", err)})
		return
	}
	claims, err := decodeToken(r)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding jwt: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding session: %s", err)})
		return
	}

	if claims.Role != model.AdminRole {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "You do not have the sufficient privileges to make a new user"})
		return
	}

	hash := HashAndSaltPassword([]byte(newUser.Password))
	log.Logger.Debug().Str("hash", hash).Msgf("Created hash for password %s", newUser.Password)
	user := model.DBUser{
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
		Username:  newUser.Username,
		Role:      newUser.Role,
	}
	err = database.AddNewUser(user, hash)
	if err != nil {
		log.Logger.Warn().Msgf("Error creating user: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error adding user to db: %s", err)})
		return
	}
	json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
}

//TODO Check up against current password, to prevent JWT hijacking
func updateUser(w http.ResponseWriter, r *http.Request) {
	updatedUser := model.APIUser{}
	hashedAndSaltedPwd := ""
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatedUser); err != nil {
		log.Logger.Warn().Msgf("Error decoding json: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding json: %s", err)})
		return
	}

	//Getting claims form jwt
	claims, err := decodeToken(r)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding jwt token: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Error parsing jwt token"})
		return
	}
	updatedUser.Username = claims.Username
	log.Logger.Debug().Msgf("Claims from jwt: %v", claims)

	log.Logger.Debug().Msgf("Username: %s, Password: %s", updatedUser.Username, updatedUser.CurrPassword)
	loginStatus, err := database.LoginUser(updatedUser.Username, updatedUser.CurrPassword)
	log.Logger.Debug().Bool("loginStatus", loginStatus).Msg("GetPasswordHash result")
	if err != nil {
		log.Logger.Warn().Msgf("Error Loggin in user: /%s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("%s", err)})
		return
	}
	if !loginStatus {
		log.Logger.Debug().Msg("Incorrect password")
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Incorrect password"})
		return
	}

	if len(updatedUser.Password) > 0 {
		log.Logger.Debug().Msg("Updating password")
		if updatedUser.Password != updatedUser.ConfirmPw {
			log.Logger.Warn().Msg("Updated passwords do not match")
			json.NewEncoder(w).Encode(model.APIResponse{Error: "Password does not match confirm password"})
			return
		}
		hashedAndSaltedPwd = HashAndSaltPassword([]byte(updatedUser.Password))
	}

	log.Logger.Debug().Msgf("Updated user is: %v", updatedUser)
	database.UpdateUser(updatedUser, hashedAndSaltedPwd)

	json.NewEncoder(w).Encode(model.APIResponse{Error: ""})

}

func deleteUser(w http.ResponseWriter,r *http.Request) {
	var user = model.APIUser{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		log.Logger.Warn().Msgf("Error decoding json: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding json: %s", err)})
		return
	}

	claims, err := decodeToken(r)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding jwt: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding session: %s")})
		return
	}
	if claims.Role != model.AdminRole {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "You do not have the sufficient privileges to delete a user"})
		return
	}

	if claims.Username == user.Username {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "You cannot remove your own user"})
		return
	}

	if err := database.RemoveUser(user.Username); err != nil {
		log.Logger.Warn().Msgf("Error deleting user from db: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error deleting user from db: %s",err)})
		return
	}

	json.NewEncoder(w).Encode(model.APIResponse{Error: ""})


}
