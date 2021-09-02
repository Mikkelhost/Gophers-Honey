package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

/*
The devices API handles everything about devices/raspberry pis

All functions should write json data to the responsewriter
*/

func devicesSubrouter(r *mux.Router){
	deviceAPI := r.PathPrefix("/api/devices").Subrouter()
	deviceAPI.HandleFunc("/getdevices", getDevices)
	deviceAPI.HandleFunc("/configure", configureDevice)
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hey du må ikke få mine fucking devices!"))
}

func configureDevice(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Du er færdig mester, ingen konfiguration til dig!"))
}