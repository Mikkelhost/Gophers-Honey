package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mikkelhost/Gophers-Honey/pkg/config"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func checkForValidIp(ipStr string) (bool, error) {
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
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func createToken(user model.DBUser) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["username"] = user.Username
	atClaims["role"] = user.Role
	atClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	atClaims["email"] = user.Email
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
	// TODO: Explain this.
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

//decodeToken Only use this function after request has passed tokenmiddleware.
func decodeToken(r *http.Request) (model.Claims, error) {
	token := extractToken(r)
	log.Logger.Debug().Str("Token", token).Msg("Decoding token")
	tokenSlice := strings.Split(token, ".")
	claims := model.Claims{}
	log.Logger.Debug().Str("Base64", tokenSlice[1]).Msg("Base64 to decode")
	claimsJson, err := base64.RawStdEncoding.DecodeString(tokenSlice[1])
	if err != nil {
		log.Logger.Warn().Msgf("Error base64 decoding: %s", err)
		return model.Claims{}, err
	}
	err = json.Unmarshal(claimsJson, &claims)
	if err != nil {
		log.Logger.Warn().Msgf("Error parsing json: %s", err)
		return model.Claims{}, err
	}
	return claims, nil
}

func verifyToken(request *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(request)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
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
		return errors.New("token invalid")
	}
	return nil
}

func tokenAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "OPTIONS" {
			return
		}
		err := tokenValid(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Error in token validity: %s", err)))
			return
		}
		next(w, r)
	}
}

// HashAndSaltPassword takes a password byte string and hashes and salts
// it using bcrypt. The hashed and salted password is returned as a string
// for storage.
func HashAndSaltPassword(pwd []byte) string {
	cost := 14
	hash, err := bcrypt.GenerateFromPassword(pwd, cost)
	if err != nil {
		log.Logger.Warn().Msgf("Error generating hash")
	}

	return string(hash)
}

// isStringInStringArray returns true if the given string appears in the
// given array. Also returns the index of the element.
func isStringInStringArray(element string, array []string) (bool, int) {
	for index, temp := range array {
		if element == temp {
			return true, index
		}
	}
	return false, 0

}

// remove takes an index and string array and removes the element at the
// index position. NB! Does not preserve order.
func remove(i int, s []string) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// addIPToWhitelist takes an IP address string as input and appends it to the
// IP whitelist in the config file. No checks on whether the IP address is
// valid so IP's should only be passed if validated first.
func addIPToWhitelist(ip string) error {
	config.Conf.IpWhitelist = append(config.Conf.IpWhitelist, ip)
	err := config.WriteConf()
	if err != nil {
		log.Logger.Warn().Msgf("Error writing to config file: %s", err)
		return err
	}
	log.Logger.Debug().Msgf("Successfully added ip: %s to IP whitelist", ip)

	return nil
}

// removeIPFromWhitelist takes an IP address string and removes it from
// the config file. No checks on whether the IP address is valid so IP's
// should only be passed if validated first.
func removeIPFromWhitelist(ip string) error {
	if result, index := isStringInStringArray(ip, config.Conf.IpWhitelist); result {
		log.Logger.Debug().Msgf("Removing IP: %s from whitelist", ip)
		remove(index, config.Conf.IpWhitelist)
		err := config.WriteConf()
		if err != nil {
			log.Logger.Warn().Msgf("Error writing to config file: %s", err)
			return err
		}
		return nil
	}
	log.Logger.Warn().Msgf("IP address not in whitelist")
	return errors.New("ip address not in whitelist")
}
