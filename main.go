package main

import (
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	"github.com/Mikkelhost/Gophers-Honey/pkg/httpserver"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
)


func main() {
	log.InitLog()
	log.Logger.Info().Msg("Setting up database")
	//Setting up database
	database.Connect()
	//database.ConfigureDb()
	defer database.Disconnect()


	log.Logger.Info().Msg("Running server")

	httpserver.RunServer()
}
