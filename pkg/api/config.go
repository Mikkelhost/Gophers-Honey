package api

import (
	"encoding/json"
	"fmt"
	"github.com/Mikkelhost/Gophers-Honey/pkg/config"
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"github.com/Mikkelhost/Gophers-Honey/pkg/notification"
	"github.com/Mikkelhost/Gophers-Honey/pkg/piimage"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func configSubrouter(r *mux.Router) {
	configAPI := r.PathPrefix("/api/config").Subrouter()
	configAPI.HandleFunc("/whitelist", whitelistHandler).Methods("PATCH", "OPTIONS")
	configAPI.HandleFunc("", configHandler).Methods("GET", "POST", "PATCH", "OPTIONS")
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// CORS preflight handling.
	if r.Method == "OPTIONS" {
		return
	}
	switch r.Method {
	case "GET":
		getConfig(w, r)
		return
	case "POST":
		setupService(w, r)
		return
	case "PATCH":
		configureSmtpServer(w, r)
		return
	}
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	log.Logger.Debug().Bool("configured", config.Conf.Configured)

	json.NewEncoder(w).Encode(model.ConfigResponse{Configured: config.Conf.Configured})
}

// setupService is a function that will be called when the
// service has been run for the first time and the user has
// has submitted all the user details and image details
func setupService(w http.ResponseWriter, r *http.Request) {
	if !config.Conf.Configured {
		var setup = model.SetupParams{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&setup); err != nil {
			log.Logger.Warn().Msgf("Failed decoding json: %s", err)
			json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Failed decoding json: %s", err)})
			return
		}
		log.Logger.Debug().Msgf("setup params: %v", setup)
		if setup.User.Password != setup.User.ConfirmPw {
			log.Logger.Warn().Msg("Passwords does not match")
			json.NewEncoder(w).Encode(model.APIResponse{Error: "Passwords need to match"})
			return
		}
		//Make first image
		port, err := strconv.Atoi(setup.Image.Port)
		if err != nil {
			log.Logger.Warn().Msgf("Error converting port to int: %s", err)
			json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("%s", err)})
			return
		}
		if port < 0 {
			log.Logger.Warn().Msg("Port smaller than 0, aborting")
			json.NewEncoder(w).Encode(model.APIResponse{Error: "Port cannot be smaller than 0"})
			return
		}
		id, err := database.NewImage(model.Image{
			Name: setup.Image.ImageName,
		})
		if err != nil {
			log.Logger.Warn().Msgf("Error creating new image in db: %s", err)
			json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("%s", err)})
			return
		}

		err = piimage.InsertConfig(model.PiConf{
			C2:        setup.Image.C2,
			Port:      port,
			DeviceID:  0,
			DeviceKey: DEVICE_KEY,
		}, id)
		if err != nil {
			log.Logger.Warn().Msgf("Error inserting config into image: %s", err)
			json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("%s", err)})
			return
		}

		//Making first user
		hash := HashAndSaltPassword([]byte(setup.User.Password))
		log.Logger.Debug().Str("hash", hash).Msgf("Created hash for password %s", setup.User.Password)
		user := model.DBUser{
			FirstName: setup.User.FirstName,
			LastName:  setup.User.LastName,
			Email:     setup.User.Email,
			Username:  setup.User.Username,
			Role:      "Admin",
		}
		err = database.AddNewUser(user, hash)
		if err != nil {
			log.Logger.Warn().Msgf("Error creating user: %s", err)
			json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("%s", err)})
			return
		}
		token, err := createToken(user)
		if err != nil {
			log.Logger.Warn().Msgf("Error creating token: ", err)
			json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("%s", err)})
			return
		}

		// Setting config
		config.Conf.Configured = true
		if err := config.WriteConf(); err != nil {
			log.Logger.Warn().Msgf("Error setting config: %s", err)
			json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("%s", err)})
			return
		}
		json.NewEncoder(w).Encode(model.APIUser{Token: token})
	} else {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Service has already been configured"})
	}
}

// configureSmtpServer sets and writes SMTP server configurations such as
// username/email, password, SMTP server and SMTP port.
func configureSmtpServer(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var smtpServer model.SmtpServer

	err := decoder.Decode(&smtpServer)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding JSON: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding JSON: %s", err)})
		return
	}

	err = notification.ConfigureSmtpServer(smtpServer)
	if err != nil {
		log.Logger.Warn().Msgf("Error configuring SMTP server: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error configuring SMTP server: %s", err)})
		return
	}

	json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
}



func whitelistHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	switch r.Method {

	case "PATCH":
		updateWhitelist(w, r)
	}
}

func updateWhitelist(w http.ResponseWriter, r *http.Request) {

}
// addIPToWhitelist takes an IP address string as input and appends it to the
// IP whitelist in the config file. No checks on whether the IP address is
// valid so IP's should only be passed if validated first.
func addIPToWhitelist(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var ip model.IPAddress

	err := decoder.Decode(&ip)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding JSON: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding JSON: %s", err)})
		return
	}

	config.Conf.IpWhitelist = append(config.Conf.IpWhitelist, ip.IPAddressString)
	err = config.WriteConf()
	if err != nil {
		log.Logger.Warn().Msgf("Error writing to config file: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error writing to config file: %s", err)})
		return
	}
	log.Logger.Debug().Msgf("Successfully added ip: %s to IP whitelist", ip)

	err = json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
}

// removeIPFromWhitelist takes an IP address string and removes it from
// the config file.
func removeIPFromWhitelist(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var ip model.IPAddress

	err := decoder.Decode(&ip)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding JSON: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding JSON: %s", err)})
		return
	}
	if result, index := isStringInStringArray(ip.IPAddressString, config.Conf.IpWhitelist); result {
		log.Logger.Debug().Msgf("Removing IP: %s from whitelist", ip)
		remove(index, config.Conf.IpWhitelist)
		err = config.WriteConf()
		if err != nil {
			log.Logger.Warn().Msgf("Error writing to config file: %s", err)
			json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error writing to config file: %s", err)})
			return
		}
		json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
		return
	}
	log.Logger.Warn().Msgf("IP address not in whitelist")
	json.NewEncoder(w).Encode(model.APIResponse{Error: "IP address not in whitelist"})
}
