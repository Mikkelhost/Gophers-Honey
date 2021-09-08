package main

import (
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	"github.com/Mikkelhost/Gophers-Honey/pkg/httpserver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)


func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Setting up database")
	//Setting up database
	database.Connect()
	//database.ConfigureDb()
	defer database.Disconnect()


	log.Info().Msg("Running server")

	httpserver.RunServer()
}