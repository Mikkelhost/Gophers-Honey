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
The logs API handles everything about logs/raspberry pis

All functions should write json data to the responsewriter
*/
func logsSubrouter(r *mux.Router) {
	logAPI := r.PathPrefix("/api/logs").Subrouter()
	logAPI.HandleFunc("/getlogs", TokenAuthMiddleware(getLogs)).Methods("GET")
	logAPI.HandleFunc("/addLog", deviceSecretMiddleware(newLog)).Methods("POST")
}

// newLog is called by devices when they create a new log
// The log message needs to be a valid JSON string
func newLog(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var logStruct database.Log

	if err := decoder.Decode(&logStruct); err != nil {
		log.Logger.Warn().Msgf("Failed decoding json: %s", err)
		w.Write([]byte(fmt.Sprintf("Failed decoding json: %s", err)))
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

func getLogs(w http.ResponseWriter, r *http.Request) {
	var logs []database.Log
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
