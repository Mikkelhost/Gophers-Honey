package main

import (
	"github.com/GeekMuch/GoHoney/pkg/database"
	"github.com/GeekMuch/GoHoney/pkg/httpserver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)




func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Setting up database")
	//Setting up database
	db := database.Connect("gohoney", "password", "127.0.0.1", "gohoney")
	database.ConfigureDb()
	defer db.Close()

	log.Info().Msg("Running server")

	httpserver.RunServer()
}
