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
The devices API handles everything about devices/raspberry pis

All functions should write json data to the responsewriter
*/

var DEVICE_KEY = getenv("DEVICE_KEY", "XxPFUhQ8R7kKhpgubt7v")

// devicesSubrouter
// Sets up a devices API subrouter
func devicesSubrouter(r *mux.Router) {
	deviceAPI := r.PathPrefix("/api/devices").Subrouter()
	deviceAPI.HandleFunc("", tokenAuthMiddleware(deviceHandler)).Methods("GET", "PUT", "DELETE", "OPTIONS")
	//deviceAPI.HandleFunc("/configure", tokenAuthMiddleware(configureDevice)).Methods("POST", "OPTIONS")
	deviceAPI.HandleFunc("/addDevice", deviceSecretMiddleware(newDevice)).Methods("POST")
	deviceAPI.HandleFunc("/getDeviceConf", deviceSecretMiddleware(getDeviceConfiguration)).Methods("GET")
	deviceAPI.HandleFunc("/heartbeat", deviceSecretMiddleware(handleHeartbeat)).Methods("POST")
	//deviceAPI.HandleFunc("/removeDevice", tokenAuthMiddleware(removeDevice)).Methods("POST")
}

func deviceHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// CORS preflight handling.
	if r.Method == "OPTIONS" {
		return
	}
	switch r.Method {
	case "GET":
		getDevices(w, r)
		return
	case "PUT":
		configureDevice(w, r)
		return
	case "DELETE":
		removeDevice(w, r)
		return
	}
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	var devices []model.Device
	devices, err := database.GetAllDevices()
	if err != nil {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Error retrieving devices"})
		return
	}
	if len(devices) == 0 {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "No devices in DB"})
		return
	}
	devicesJson, err := json.Marshal(devices)
	if err != nil {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Error Marshalling devices"})
		return
	}

	w.Write(devicesJson)
}

// getDeviceConfiguration retrieves the configuration information stored
// for a specific device and sends a JSON response containing the device
// configuration to the requester.
func getDeviceConfiguration(w http.ResponseWriter, r *http.Request) {
	var configuration model.Configuration
	var device = model.DeviceAuth{}
	var err error

	decoder := json.NewDecoder(r.Body)

	if err = decoder.Decode(&device); err != nil {
		log.Logger.Warn().Msgf("Error decoding JSON: %s", err)
		return
	}

	configuration, err = database.GetDeviceConfiguration(device.DeviceId)

	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}

	response := model.PiConfResponse{
		Status:   "Success",
		DeviceId: device.DeviceId,
		Services: model.Service{
			SSH:    configuration.Services.SSH,
			FTP:    configuration.Services.FTP,
			SMB:    configuration.Services.SMB,
			RDP:    configuration.Services.RDP,
			TELNET: configuration.Services.TELNET,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// configureDevice updates the configured services for a specified
// device ID.
func configureDevice(w http.ResponseWriter, r *http.Request) {
	var config model.Configuration

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&config); err != nil {
		log.Logger.Warn().Msgf("Error decoding JSON: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding JSON: %s", err)})
		return
	}

	err := database.ConfigureDevice(config.Services, config.DeviceID)
	if err != nil {
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error updating device configuration: %s", err)})
		return
	}

	json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
}

// newDevice
// Called by devices when they are booted up for the first time
// Will return an unique id for the device to use going forward
func newDevice(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var ipStruct = model.Device{}
	if err := decoder.Decode(&ipStruct); err != nil {
		log.Logger.Warn().Msgf("Error decoding JSON: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding JSON: %s", err)})
		return
	}

	//Checking if the ipstr is a valid IP address
	found, err := checkForValidIp(ipStruct.IpStr)
	if err != nil {
		log.Logger.Warn().Msgf("Error in regex: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Internal server error"})
		return
	}
	if !found {
		log.Logger.Debug().Msg("Ip is invalid")
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Ip is invalid"})
		return
	}

	log.Logger.Debug().Msgf("Received new device with IP: %s Adding to DB", ipStruct.IpStr)
	deviceID, err := database.AddDevice(strings.TrimSpace(ipStruct.IpStr))
	if err != nil {
		log.Logger.Warn().Msgf("Error adding device: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Error adding device"})
		return
	}
	response := model.PiConfResponse{
		Status:   "Success",
		DeviceId: deviceID,
	}
	json.NewEncoder(w).Encode(response)
}

// removeDevice
// TODO: in progress
func removeDevice(w http.ResponseWriter, r *http.Request) {
	var deviceID uint32
	var device = model.Device{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&device); err != nil {
		log.Logger.Warn().Msgf("Error decoding JSON: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding JSON: %s", err)})
		return
	}
	log.Logger.Debug().Uint32("device_id", device.DeviceID).Msg("Deleting device")
	deviceID = device.DeviceID
	err := database.RemoveDevice(deviceID)
	if err != nil {
		log.Logger.Warn().Msgf("Error removing device: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Internal server error"})
		return
	}

	json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
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
			w.Write([]byte("Please provide the right credentials"))
			return
		}
		next(w, r)
	}
}

// handleHeartbeat handles device heartbeats for a given device and calls
// the database handler to update the "last_seen" timestamp for the device.
func handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	var heartbeat model.Heartbeat

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&heartbeat); err != nil {
		log.Logger.Warn().Msgf("Error decoding JSON: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding JSON: %s", err)})
		return
	}

	log.Logger.Debug().Uint32("device_id", heartbeat.DeviceID).Msg("Received heartbeat from device")

	err := database.HandleHeartbeat(heartbeat.DeviceID)
	if err != nil {
		log.Logger.Warn().Msgf("Error handling heartbeat: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error handling heartbeat: %s", err)})
		return
	}

	log.Logger.Debug().Msg("Heartbeat successfully handled")
	ClientPool.Heartbeat <- heartbeat.DeviceID
	json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
}
