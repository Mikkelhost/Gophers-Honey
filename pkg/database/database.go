package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)
var db *sql.DB

func Connect(username, password, host, dbName string) *sql.DB{
	log.Info().Msgf("Connecting to sql server with: %s, %s, %s, %s", username, password, host, dbName)
	connectString := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, dbName)
	log.Info().Msgf("Connectstring: %s", connectString)
	database, err := sql.Open("mysql", connectString)
	if err != nil {
		log.Fatal().Msgf("Error: %s", err)
	}
	db = database
	return db
}

func ConfigureDb() {
	log.Info().Msg("Configuring DB")
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS devices (id INT NOT NULL PRIMARY KEY, configured BOOL, services CHAR(255), ipAddress CHAR(255));"); err != nil {
		log.Fatal().Msgf("Error: %s", err)
	}
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS backend (configured BOOL);"); err != nil {
		log.Fatal().Msgf("Error: %s", err)
	}

	var conf bool

	row := db.QueryRow("SELECT * FROM backend")
	if err := row.Scan(&conf); err == sql.ErrNoRows || err != nil {
		log.Fatal().Msgf("Error: %s", err)
		log.Info().Msg("No config in database, inserting not configured to DB")
		if err == sql.ErrNoRows {
			if _, err := db.Exec("INSERT INTO backend (configured) VALUES(false)"); err != nil {
				log.Fatal().Msgf("Error: %s", err)
			}
		}
	}
}