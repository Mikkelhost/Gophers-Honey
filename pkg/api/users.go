package api

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

/*
The users API handles everything about users

All functions should write json data to the responsewriter
*/

type UserAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

	w.Write([]byte("fucker"))
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Registering user"))
}

func createToken(userid uint32) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}
	return token, nil
}

func extractToken(request *http.Request) string {
	bearToken := request.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(request *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(request)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func tokenValid(request *http.Request) error {
	token, err := verifyToken(request)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return errors.New("Token invalid")
	}
	return nil
}

func tokenAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tokenValid(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Error in token validity: %s", err)))
			return
		}
		next(w, r)
	}
}
