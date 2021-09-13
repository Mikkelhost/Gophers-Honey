package api

import (
	"encoding/json"
	"fmt"
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

/*
The users API handles everything about users

All functions should write json data to the responsewriter
*/

type UserAuth struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token string `json:"token,omitempty"`
	Error string  `json:"error"`
}

func usersSubrouter(r *mux.Router) {
	usersAPI := r.PathPrefix("/api/users").Subrouter()
	usersAPI.HandleFunc("/getUsers", getUsers).Methods("GET")
	usersAPI.HandleFunc("/login", loginUser).Methods("POST", "OPTIONS")
	usersAPI.HandleFunc("/register", registerUser).Methods("POST")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Getting users"))
}

// TODO Finish the loginUser api endpoint, return token to user.
func loginUser(w http.ResponseWriter, r *http.Request) {
	var userInfo = UserAuth{}
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
	hash, err := database.GetPasswordHash(strings.TrimSpace(userInfo.Username))
	log.Logger.Debug().Str("hash",hash).Msg("GetPasswordHash result")
	if err != nil {
		log.Logger.Warn().Msgf("Error getting password hash: /%s", err)
		//TODO return something to user
		json.NewEncoder(w).Encode(UserAuth{Error: "Username or Password does not exists"})
		return
	}
	w.Write([]byte("fucker"))
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Registering user"))
}
