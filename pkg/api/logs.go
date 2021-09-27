package api

import (
	"encoding/json"
	"fmt"
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

/*
The logs API handles everything about logs/raspberry pis

All functions should write json data to the responseWriter
*/
func logsSubrouter(r *mux.Router) {
	logAPI := r.PathPrefix("/api/logs").Subrouter()
	logAPI.HandleFunc("/getlogs", tokenAuthMiddleware(getLogs)).Methods("GET")
	logAPI.HandleFunc("/updateTTLIndex", tokenAuthMiddleware(updateTTLIndex)).Methods("POST")
	logAPI.HandleFunc("/addLog", deviceSecretMiddleware(newLog)).Methods("POST")
}

// newLog is called by devices when they create a new log
// The log message needs to be a valid JSON string
func newLog(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var logStruct model.Log

	if err := decoder.Decode(&logStruct); err != nil {
		log.Logger.Warn().Msgf("Error decoding json: %s", err)
		w.Write([]byte(fmt.Sprintf("Error decoding json: %s", err)))
		return
	}

	log.Logger.Debug().Msgf("Received new log for device ID: %s Adding to DB", logStruct.DeviceID)
	err := database.AddLog(logStruct.DeviceID, logStruct.TimeStamp, strings.TrimSpace(logStruct.Message))
	if err != nil {
		log.Logger.Warn().Msgf("Error adding log: %s", err)
		w.Write([]byte("Error adding log"))
		return
	}

	w.Write([]byte(fmt.Sprintf("{\"status\": \"Success\", \"device_id\": %d,\"log_id\": %d }",
		logStruct.DeviceID, logStruct.LogID)))

}

// getLogs retrieves all logs currently present in the database.
func getLogs(w http.ResponseWriter, r *http.Request) {
	var logs []model.Log
	logs, err := database.GetAllLogs()
	if err != nil {
		w.Write([]byte("Error retrieving devices"))
		return
	}
	if len(logs) == 0 {
		w.Write([]byte("No logs in DB"))
		return
	}
	logsJson, err := json.Marshal(logs)
	if err != nil {
		w.Write([]byte("Error Marshalling logs"))
		return
	}

	w.Write(logsJson)
}

// updateTTLIndex updates the "setExpireAfterSeconds" index of the
//// "log_collection" collection
func updateTTLIndex(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var update struct {
		ExpireAfterSeconds int32 `json:"expire_after_seconds"`
	}

	err := decoder.Decode(&update)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding json: %s", err)
		w.Write([]byte(fmt.Sprintf("Error decoding json: %s", err)))
		return
	}

	log.Logger.Debug().Int32("expireAfterSeconds", update.ExpireAfterSeconds).Msgf("Value decoded as:")

	err = database.UpdateTTLIndex(update.ExpireAfterSeconds)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding json: %s", err)
		w.Write([]byte(fmt.Sprintf("Error decoding json: %s", err)))
		return
	}
}
