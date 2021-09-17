package api

import (
	"encoding/json"
	"fmt"
	"github.com/Mikkelhost/Gophers-Honey/pkg/config"
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
)

type ConfigResponse struct {
	Configured bool `json:"configured"`
}

type SetupParams struct {
	Image ImageInfo `json:"imageInfo"`
	User  User  `json:"userInfo"`
}

type ImageInfo struct {
	ImageName string `json:"name"`
	Hostname  string `json:"hostname"`
	Port      string `json:"port"`
}

func configSubrouter(r *mux.Router) {
	configAPI := r.PathPrefix("/api/config").Subrouter()
	configAPI.HandleFunc("/getConfig", getConfig).Methods("GET", "OPTIONS")
	configAPI.HandleFunc("/setupService", setupService).Methods("POST", "OPTIONS")
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	log.Logger.Debug().Bool("configured", config.Conf.Configured)

	json.NewEncoder(w).Encode(ConfigResponse{Configured: config.Conf.Configured})
}

// setupService is a function that will be called when the
// service has been run for the first time and the user has
// has submitted all the user details and image details
func setupService(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	if !config.Conf.Configured {
		var setup = SetupParams{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&setup); err != nil {
			log.Logger.Warn().Msgf("Failed decoding json: %s", err)
			w.Write([]byte(fmt.Sprintf("Failed decoding json: %s", err)))
			return
		}
		log.Logger.Debug().Msgf("setup params: %v", setup)

		//Making first user
		hash := HashAndSaltPassword([]byte(setup.User.Password))
		log.Logger.Debug().Str("hash", hash).Msgf("Created hash for password %s", setup.User.Password)
		err := database.AddNewUser(database.User{
			FirstName: setup.User.FirstName,
			LastName: setup.User.LastName,
			Email: setup.User.Email,
			Username: setup.User.Username,
		}, hash)
		if err != nil {
			log.Logger.Warn().Msgf("Error creating user: %s", err)
			json.NewEncoder(w).Encode(User{Error: fmt.Sprintf("%s", err)})
			return
		}
		token, err := createToken(setup.User.Username)
		if err != nil {
			log.Logger.Warn().Msgf("Error creating token: ", err)
			json.NewEncoder(w).Encode(User{Error: fmt.Sprintf("%s", err)})
			return
		}

		//Make first image


		//Setting config
		conff := config.Config{
			Configured: true,
		}
		if err := config.SetConfig(conff); err != nil {
			log.Logger.Warn().Msgf("Error setting config: %s", err)
			json.NewEncoder(w).Encode(User{Error: fmt.Sprintf("%s",err)})
			return
		}
		json.NewEncoder(w).Encode(User{Token: token})
	} else {
		json.NewEncoder(w).Encode(User{Error: "Service has already been configured"})
	}
}
