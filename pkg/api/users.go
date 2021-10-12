package api

import (
	"encoding/json"
	"fmt"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"net/http"

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
	usersAPI.HandleFunc("",tokenAuthMiddleware(userHandler)).Methods("GET", "POST", "OPTIONS")
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
	}
}
// getUsers gets all information of all users, but the hashed passwords
func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []model.DBUser
	users, err := database.GetAllUsers()
	if err != nil {
		w.Write([]byte("Error retrieving users"))
		return
	}
	if len(users) == 0 {
		w.Write([]byte("No users in DB"))
		return
	}
	usersJson, err := json.Marshal(users)
	if err != nil {
		w.Write([]byte("Error Marshalling users"))
		return
	}
	w.Write(usersJson)
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
		w.Write([]byte(fmt.Sprintf("Failed decoding json: %s", err)))
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

	token, err := createToken(userInfo.Username)
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
	var newUser = model.DBUser{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser); err != nil {
		log.Logger.Warn().Msgf("Error decoding json: %s", err)
		w.Write([]byte(fmt.Sprintf("Error decoding json: %s", err)))
		return
	}

	hashedAndSaltedPassword := HashAndSaltPassword([]byte(newUser.PasswordHash))

	err := database.AddNewUser(newUser, hashedAndSaltedPassword)
	if err != nil {
		log.Logger.Warn().Msgf("Error registering user: %s", err)
		w.Write([]byte(fmt.Sprintf("Error registering user: %s", err)))
	}
	w.Write([]byte("Registering user"))
}
