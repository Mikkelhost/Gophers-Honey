package api

import (
	"encoding/json"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"

)

type ConfigResponse struct {
	Configured bool `json:"configured"`
}

func configSubrouter(r *mux.Router) {
	configAPI := r.PathPrefix("/api/config").Subrouter()
	configAPI.HandleFunc("/getConfig", getConfig).Methods("GET", "OPTIONS")
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	log.Logger.Debug().Bool("configured", conf.Configured)

	json.NewEncoder(w).Encode(ConfigResponse{Configured: conf.Configured})
}