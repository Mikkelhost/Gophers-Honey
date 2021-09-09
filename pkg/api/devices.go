package api

import (
	"encoding/json"
	"fmt"
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	"github.com/gorilla/mux"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"net/http"
)
/*
The devices API handles everything about devices/raspberry pis

All functions should write json data to the responsewriter
*/
type DeviceAuth struct {
	DeviceId uint32 `json:"device_id"`
	DeviceKey string `json:"deviceKey"`
}
var DEVICE_KEY = getenv("DEVICE_KEY","XxPFUhQ8R7kKhpgubt7v")

func devicesSubrouter(r *mux.Router){
	deviceAPI := r.PathPrefix("/api/devices").Subrouter()
	deviceAPI.HandleFunc("/getdevices", tokenAuthMiddleware(getDevices)).Methods("GET")
	deviceAPI.HandleFunc("/configure", tokenAuthMiddleware(configureDevice)).Methods("POST")
	deviceAPI.HandleFunc("/addDevice", deviceSecretMiddleware(newDevice)).Methods("POST")
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hey du må ikke få mine fucking devices!"))
}

func configureDevice(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Du er færdig mester, ingen konfiguration til dig!"))
}

// TODO Make sure that ip is not empty.
func newDevice(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var ipStruct = database.Device{}
	if err := decoder.Decode(&ipStruct); err != nil {
		log.Logger.Warn().Msgf("Failed decoding json: %s", err)
		w.Write([]byte(fmt.Sprintf("Failed decoding json: %s", err)))
		return
	}
	w.Write([]byte(fmt.Sprintf("Success")))
	log.Logger.Debug().Msgf("Received new device with IP: %s Adding to DB", ipStruct.IpStr)
	database.AddDevice(ipStruct.IpStr)
}

func deviceSecretMiddleware(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		deviceKey := extractToken(r)
		log.Logger.Debug().Msgf("Received authentication attempt with key: %s", deviceKey)
		if deviceKey != DEVICE_KEY {
			log.Logger.Debug().Msg("Wrong Devicekey for authentication")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Please provide the right credential"))
			return
		}
		next(w, r)
	}
}