package api

import (
	"encoding/json"
	"fmt"
	"github.com/Mikkelhost/Gophers-Honey/pkg/config"
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"github.com/Mikkelhost/Gophers-Honey/pkg/notification"
	"github.com/gorilla/mux"
	"net/http"
)

/*
The logs API handles everything about logs/raspberry pis

All functions should write json data to the responseWriter
*/
func logsSubrouter(r *mux.Router) {
	logAPI := r.PathPrefix("/api/logs").Subrouter()
	logAPI.HandleFunc("", tokenAuthMiddleware(logHandler)).Methods("GET", "PUT", "OPTIONS")
	logAPI.HandleFunc("/addLog", deviceSecretMiddleware(newLog)).Methods("POST")
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// CORS preflight handling.
	if r.Method == "OPTIONS" {
		return
	}
	switch r.Method {
	case "GET":
		getLogs(w, r)
		return
	case "PUT":
		updateTTLIndex(w, r)
		return
	}
}

// newLog is called by devices when they create a new log
// The log message needs to be a valid JSON string
func newLog(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var newLog model.Log

	if err := decoder.Decode(&newLog); err != nil {
		log.Logger.Warn().Msgf("Error decoding json: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding json: %s", err)})
		return
	}

	log.Logger.Debug().Msgf("Received new log from device ID: %s Adding to DB", newLog.DeviceID)
	err := database.AddLog(newLog)
	if err != nil {
		log.Logger.Warn().Msgf("Error adding log: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Error adding log"})
		return
	}

	if newLog.Level == model.CRITICAL {
		log.Logger.Info().Msgf("Critical level log received. Notifying users.")
		err = notification.NotifyAll(newLog)
		if err != nil {
			log.Logger.Warn().Msgf("Error notifying users: %s", err)
			json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding JSON: %s", err)})
			return
		}
	} else if newLog.Level == model.SCAN {
		// Send no alert if source ip appears in the whitelist.
		if result, _ := isStringInStringArray(newLog.SrcHost, config.Conf.IpWhitelist); result {
			err = json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
			if err != nil {
				log.Logger.Warn().Msgf("Error encoding json: %s", err)
				return
			}
		} else {
			log.Logger.Info().Msgf("Scan level log received. Notifying users")
			err := notification.NotifyAll(newLog)
			if err != nil {
				log.Logger.Warn().Msgf("Error notifying users: %s", err)
				json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding JSON: %s", err)})
				return
			}
		}
	}

	json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
}

// getLogs retrieves all logs currently present in the database.
func getLogs(w http.ResponseWriter, r *http.Request) {
	var logs []model.Log
	logs, err := database.GetAllLogs()
	if err != nil {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Error retrieving logs"})
		return
	}
	if len(logs) == 0 {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "No logs in DB"})
		return
	}

	json.NewEncoder(w).Encode(logs)
}

// updateTTLIndex updates the "setExpireAfterSeconds" index of the
// "log_collection" collection
func updateTTLIndex(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var update struct {
		ExpireAfterSeconds int32 `json:"expire_after_seconds"`
	}

	err := decoder.Decode(&update)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding json: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding json: %s", err)})
		return
	}

	log.Logger.Debug().Int32("expireAfterSeconds", update.ExpireAfterSeconds).Msgf("Value decoded as:")

	err = database.UpdateTTLIndex(update.ExpireAfterSeconds)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding json: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding json: %s", err)})
		return
	}

	json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
}
