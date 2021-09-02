package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

/*
The users API handles everything about users

All functions should write json data to the responsewriter
 */

func usersSubrouter(r *mux.Router) {
	usersAPI := r.PathPrefix("/api/users").Subrouter()
	usersAPI.HandleFunc("/getUsers", getUsers)
	usersAPI.HandleFunc("/login", loginUser)
	usersAPI.HandleFunc("/register", registerUser)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Getting users"))
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logging in user"))
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Registering user"))
}
