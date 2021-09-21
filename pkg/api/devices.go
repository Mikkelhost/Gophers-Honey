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
The devices API handles everything about devices/raspberry pis

All functions should write json data to the responsewriter
*/

type DeviceAuth struct {
	DeviceId  uint32 `json:"device_id"`
	DeviceKey string `json:"deviceKey"`
}

var DEVICE_KEY = getenv("DEVICE_KEY", "XxPFUhQ8R7kKhpgubt7v")

// devicesSubrouter
// Sets up a devices API subrouter
func devicesSubrouter(r *mux.Router) {
	deviceAPI := r.PathPrefix("/api/devices").Subrouter()
	deviceAPI.HandleFunc("/getDevices", tokenAuthMiddleware(getDevices)).Methods("GET", "OPTIONS")
	deviceAPI.HandleFunc("/configure", tokenAuthMiddleware(configureDevice)).Methods("POST", "OPTIONS")
	deviceAPI.HandleFunc("/addDevice", deviceSecretMiddleware(newDevice)).Methods("POST")
	deviceAPI.HandleFunc("/getDeviceConf", deviceSecretMiddleware(getDeviceConfiguration)).Methods("POST")
	deviceAPI.HandleFunc("/removeDevice", tokenAuthMiddleware(removeDevice)).Methods("POST")
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	var devices []database.Device
	devices, err := database.GetAllDevices()
	if err != nil {
		w.Write([]byte("Error retrieving devices"))
		return
	}
	if len(devices) == 0 {
		w.Write([]byte("No devices in DB"))
		return
	}
	devicesJson, err := json.Marshal(devices)
	if err != nil {
		w.Write([]byte("Error Marshalling devices"))
		return
	}

	w.Write(devicesJson)
}

// getDeviceConfiguration retrieves the configuration information stored
// for a specific device and sends a JSON response containing the device
// configuration to the requester.
func getDeviceConfiguration(w http.ResponseWriter, r *http.Request) {
	var configuration database.Configuration
	var deviceID uint32
	var err error

	decoder := json.NewDecoder(r.Body)

	if err = decoder.Decode(&deviceID); err != nil {
		log.Logger.Warn().Msgf("Error decoding JSON: %s", err)
	}

	configuration, err = database.GetDeviceConfiguration(deviceID)

	if err != nil {
		return
	}

	w.Write([]byte(fmt.Sprintf("{\"status\": \"Success\", \"device_id\": %d, \"services\": {"+
		"\"SSH:\": %t, "+
		"\"FTP:\": %t, "+
		"\"Telnet:\": %t, "+
		"\"RDP:\": %t, "+
		"\"SMB:\": %t }",
		deviceID, configuration.Services.SSH, configuration.Services.FTP,
		configuration.Services.TELNET, configuration.Services.RDP, configuration.Services.SMB)))
}

func configureDevice(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	w.Write([]byte("Du er færdig mester, ingen konfiguration til dig!"))
}

// newDevice
// Called by devices when they are booted up for the first time
// Will return an unique id for the device to use going forward
func newDevice(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var ipStruct = database.Device{}
	if err := decoder.Decode(&ipStruct); err != nil {
		log.Logger.Warn().Msgf("Error decoding JSON: %s", err)
		w.Write([]byte(fmt.Sprintf("Error decoding JSON: %s", err)))
		return
	}

	//Checking if the ipstr is a valid IP address
	found, err := checkForValidIp(ipStruct.IpStr)
	if err != nil {
		log.Logger.Warn().Msgf("Error in regex: %s", err)
		w.Write([]byte("Internal server error"))
		return
	}
	if !found {
		log.Logger.Debug().Msg("Ip is invalid")
		w.Write([]byte("Ip is invalid"))
		return
	}

	log.Logger.Debug().Msgf("Received new device with IP: %s Adding to DB", ipStruct.IpStr)
	deviceID, err := database.AddDevice(strings.TrimSpace(ipStruct.IpStr))
	if err != nil {
		log.Logger.Warn().Msgf("Error adding device: %s", err)
		w.Write([]byte("Error adding device"))
		return
	}

	w.Write([]byte(fmt.Sprintf("{\"status\": \"Success\", \"device_id\": %d}", deviceID)))
}

// removeDevice
// TODO: in progress
func removeDevice(w http.ResponseWriter, r *http.Request) {
	var deviceID uint32
	decoder := json.NewDecoder(r.Body)
	var ipStruct = database.Device{}
	if err := decoder.Decode(&ipStruct); err != nil {
		log.Logger.Warn().Msgf("Error decoding JSON: %s", err)
		w.Write([]byte(fmt.Sprintf("Error decoding JSON: %s", err)))
		return
	}
	deviceID = ipStruct.DeviceID
	err := database.RemoveDevice(deviceID)
	if err != nil {
		log.Logger.Warn().Msgf("Error removing device: %s", err)
		w.Write([]byte("Internal server error"))
		return
	}
	w.Write([]byte(fmt.Sprintf("{\"status\": \"Success\", \"device_id\": \"%d removed\"}", deviceID)))
}

// deviceSecretMiddleware
// Middleware function for authenticating devices before they get access
// to the end API call.
func deviceSecretMiddleware(next http.HandlerFunc) http.HandlerFunc {
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
