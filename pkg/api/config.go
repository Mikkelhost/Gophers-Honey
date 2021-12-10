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

//configSubrouter
//Routes the config api endpoints to their respective handlers.
func configSubrouter(r *mux.Router) {
	configAPI := r.PathPrefix("/api/config").Subrouter()
	configAPI.HandleFunc("/whitelist", tokenAuthMiddleware(whitelistHandler)).Methods("PATCH", "OPTIONS")
	configAPI.HandleFunc("/testEmail", tokenAuthMiddleware(testEmailHandler)).Methods("GET", "OPTIONS")
	configAPI.HandleFunc("", tokenAuthMiddleware(configHandler)).Methods("GET", "POST", "PATCH", "OPTIONS")
	configAPI.HandleFunc("/configured", getConfigured).Methods("GET", "OPTIONS")
}

//configHandler
//Handles config related REST functionalities
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

//getConfig
//Serves the whole server config over the api
func getConfig(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(model.Config{
		Configured: config.Conf.Configured,
		SmtpServer: model.SmtpServer{
			SmtpHost: config.Conf.SmtpServer.SmtpHost,
			SmtpPort: config.Conf.SmtpServer.SmtpPort,
		},
		IpWhitelist: config.Conf.IpWhitelist,
	})
}

//getConfigured
//Serves only the configured variable over the api.
func getConfigured(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	log.Logger.Debug().Bool("configured", config.Conf.Configured).Msg("Configured is")
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
	claims, err := decodeToken(r)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding jwt token: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Error decoding jwt token"})
		return
	}
	if claims.Role != model.AdminRole {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "You do not have the sufficient privileges to make this request"})
		return
	}
	decoder := json.NewDecoder(r.Body)

	var smtpServer model.SmtpServer

	err = decoder.Decode(&smtpServer)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding JSON: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding JSON: %s", err)})
		return
	}

	if len(smtpServer.SmtpHost) == 0 {
		smtpServer.SmtpHost = config.Conf.SmtpServer.SmtpHost
	}

	if len(smtpServer.Password) == 0 {
		smtpServer.Password = config.Conf.SmtpServer.Password
	}

	if len(smtpServer.Username) == 0 {
		smtpServer.Username = config.Conf.SmtpServer.Username
	}

	log.Logger.Debug().Msgf("Sending updated smtpserver conf to file: %v", smtpServer)

	err = notification.ConfigureSmtpServer(smtpServer)
	if err != nil {
		log.Logger.Warn().Msgf("Error configuring SMTP server: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error configuring SMTP server: %s", err)})
		return
	}

	json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
}

//testEmailHandler
//Sends a test email to test if the email config has been set correctly
func testEmailHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// CORS preflight handling.
	if r.Method == "OPTIONS" {
		return
	}

	claims, err := decodeToken(r)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding jwt token: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Error decoding jwt token"})
		return
	}
	log.Logger.Info().Str("Email", claims.Email).Msg("Sending test email")
	var to []string
	to = append(to, claims.Email)
	err = notification.SendTestEmail(to)
	if err != nil {
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error sending email: %s", err)})
		return
	}

	json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
}

//whitelistHandler
//Handles incoming whitelist related requests
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

// updateWhitelist adds or removes IP addresses from the whitelist based
// on the whether "delete" field is set. Returns an error if IP is not
// valid.
func updateWhitelist(w http.ResponseWriter, r *http.Request) {
	claims, err := decodeToken(r)
	if err != nil {
		log.Logger.Warn().Msgf("Error decoding jwt token: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: "Error decoding jwt token"})
		return
	}
	if claims.Role != model.AdminRole {
		json.NewEncoder(w).Encode(model.APIResponse{Error: "You do not have the sufficient privileges to make this request"})
		return
	}

	decoder := json.NewDecoder(r.Body)

	var ip model.IPAddress
  
	err = decoder.Decode(&ip)

	if err != nil {
		log.Logger.Warn().Msgf("Error decoding JSON: %s", err)
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error decoding JSON: %s", err)})
		return
	}
	validIP, err := checkForValidIp(ip.IPAddressString)
	if err != nil {
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error when checking ip against regex: %s", err)})
		return
	}

	if validIP {
		switch ip.Delete {
		case false:
			err = addIPToWhitelist(ip.IPAddressString)
			if err != nil {
				log.Logger.Warn().Msgf("Error adding IP to whitelist: %s", err)
				json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error adding IP to whitelist: %s", err)})
				return
			}
			json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
		case true:
			err = removeIPFromWhitelist(ip.IPAddressString)
			if err != nil {
				log.Logger.Warn().Msgf("Error removing IP from whitelist: %s", err)
				json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("Error removing IP from whitelist: %s", err)})
				return
			}
			json.NewEncoder(w).Encode(model.APIResponse{Error: ""})
		}
	} else {
		json.NewEncoder(w).Encode(model.APIResponse{Error: fmt.Sprintf("IP address not valid")})
		return
	}

}
