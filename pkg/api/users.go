package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/gorilla/mux"
)

/*
The users API handles everything about users

All functions should write json data to the responsewriter
*/

type User struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	ConfirmPw string `json:"confirmPw,omitempty"`
	Token     string `json:"token,omitempty"`
	Error     string `json:"error"`
}

func usersSubrouter(r *mux.Router) {
	usersAPI := r.PathPrefix("/api/users").Subrouter()
	usersAPI.HandleFunc("/getUsers", getUsers).Methods("GET")
	usersAPI.HandleFunc("/login", loginUser).Methods("POST", "OPTIONS")
	usersAPI.HandleFunc("/register", registerUser).Methods("POST")
}

// getUsers gets all information of all users, but the hashed passwords
func getUsers(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

	var users []database.User
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

// TODO Finish the loginUser api endpoint, return token to user.
// loginUser Authenticates a user by their username and password.
// It returns a JWT token for the user to use as a session cookie.
// The JWT will be used forward as authentication for authenticated API endpoints.
func loginUser(w http.ResponseWriter, r *http.Request) {
	var userInfo = User{}
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
		//TODO return something to user
		json.NewEncoder(w).Encode(User{Error: fmt.Sprintf("%s", err)})
		return
	}
	if !loginStatus {
		log.Logger.Debug().Msg("Incorrect username or password")
		json.NewEncoder(w).Encode(User{Error: "Incorrect username or password"})
		return
	}

	token, err := createToken(userInfo.Username)
	if err != nil {
		log.Logger.Warn().Msgf("Error creating token: %s", err)
		json.NewEncoder(w).Encode(User{Error: fmt.Sprintf("%s", err)})
		return
	}
	json.NewEncoder(w).Encode(User{Token: token})
}

// registerUser creates a new user from the given information,
// and hashes and salts password.
func registerUser(w http.ResponseWriter, r *http.Request) {

	var newUser = database.User{}
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}

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
