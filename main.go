package main

import (
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
)

var DEBUG = true

func main() {
	// Initialize logger and set logging level.
	log.InitLog(DEBUG)

	// Set up database connection.
	log.Logger.Info().Msg("Setting up database connection")
	database.Connect()
	defer database.Disconnect()
	database.Test()

	//log.Logger.Info().Msg("Running server")
	//httpserver.RunServer()
}
