package api

import (
	"errors"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
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