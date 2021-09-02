package database

import "github.com/rs/zerolog/log"

func GetConf() bool{
	var conf bool
	row := db.QueryRow("SELECT * FROM backend")
	if err := row.Scan(&conf); err != nil {
		log.Warn().Msgf("Error: %s", err)
	}
	return conf
}