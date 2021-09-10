package api

import (
	"errors"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"net/http"
	"regexp"
	"strings"
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