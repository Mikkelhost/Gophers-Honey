package api

import (
	"errors"
	"fmt"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func checkForValidIp(ipStr string) (bool, error)  {
	if strings.TrimSpace(ipStr) == "" {
		return false, errors.New("Error: Empty IP")
	}
	regex := "^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"
	log.Logger.Debug().Msgf("Checking if ip matches regex: %s", strings.TrimSpace(ipStr))
	found, err := regexp.Match(regex, []byte(strings.TrimSpace(ipStr)))
	log.Logger.Debug().Bool("found", found).Msg("Found is")
	if err != nil {
		return false, err
	}
	return found, nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
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

// hashAndSaltPassword takes a password byte string and hashes and salts
// it using bcrypt. The hashed and salted password is returned as a string
// for storage.
func hashAndSaltPassword(pwd []byte) string {
	cost := 14
	hash, err := bcrypt.GenerateFromPassword(pwd, cost)
	if err != nil {
		log.Logger.Warn().Msgf("Error generating hash")
	}

	return string(hash)
}