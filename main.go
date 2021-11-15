package main

import (
	"github.com/Mikkelhost/Gophers-Honey/pkg/config"
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	"github.com/Mikkelhost/Gophers-Honey/pkg/hpfeedsHandler"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"time"
)

var DEBUG = true

func main() {
	// Create config file.
	err := config.CreateConfFile()
	if err != nil {
		return
	}

	// Initialize logger and set logging level.
	log.InitLog(DEBUG)

	// Set up database connection.
	log.Logger.Info().Msg("Setting up database connection")
	database.Connect()
	defer database.Disconnect()

	// Set up server.
	log.Logger.Info().Msg("Running server")
	c, err := config.GetServiceConfig()

	if err != nil {
		log.Logger.Fatal().Msgf("Error getting config: %s", err)
	}
	if !c.Configured {
		log.Logger.Info().Msg("Service has not yet been configured, access the webpage and follow " +
			"the setup")
	}

	go func() {
		err := hpfeedsHandler.Broker()
		if err != nil {
			log.Logger.Fatal().Msgf("Broker error: %s", err)
		}
	}()
	time.Sleep(5 * time.Second)

	err = hpfeedsHandler.TestSubscriber()
	if err != nil {
		log.Logger.Fatal().Msgf("Subscriber error: %s", err)
	}

	// httpserver.RunServer()
}
