package config

import (
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var Conf *model.Config

//CreateConfFile
//Creates a config file for the backend
func CreateConfFile() error {
	// Checking if file already exists
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		f, err := os.Create("config.yml")
		if err != nil {
			log.Logger.Fatal().Msgf("Error creating config.yml: %s", err)
		}
		var config = model.Config{
			Configured: false,
		}
		yml, err := yaml.Marshal(&config)
		if err != nil {
			log.Logger.Fatal().Msgf("Error marshalling yaml string: %s", err)
		}
		_, err = f.Write(yml)
		if err != nil {
			log.Logger.Fatal().Msgf("Error writing to file: %s", err)
			return err
		}

		err = f.Close()
		if err != nil {
			log.Logger.Fatal().Msgf("Error closing file: %s", err)
			return err
		}
		return nil
	}
	return nil
}

//GetServiceConfig
//Reads the config file and assigns the global config variable to the result
func GetServiceConfig() (*model.Config, error) {
	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Logger.Warn().Msgf("Error reading file: %s", err)
		return nil, err
	}

	var config model.Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Logger.Warn().Msgf("Error unmarshalling yaml: %s", err)
		return nil, err
	}
	Conf = &config
	return &config, nil
}

//WriteConf
//Writes new config parameters to the config file.
func WriteConf() error {
	log.Logger.Debug().Msgf("Config to be set: %v", *Conf)
	log.Logger.Debug().Msgf("conf: %falsev", *Conf)
	f, err := os.OpenFile("config.yml", os.O_RDWR, 0644)
	if err != nil {
		log.Logger.Warn().Msgf("Error opening config.yml: %s", err)
		return err
	}
	yml, err := yaml.Marshal(Conf)
	if err != nil {
		log.Logger.Warn().Msgf("Error marshalling yaml string: %s", err)
		return err
	}

	err = f.Truncate(0)
	if err != nil {
		log.Logger.Warn().Msgf("Error truncating file: %s", err)
		return err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		log.Logger.Warn().Msgf("Error seeking in file: %s", err)
		return err
	}

	if _, err := f.Write(yml); err != nil {
		log.Logger.Warn().Msgf("Error writing conf", err)
		return err
	}

	err = f.Close()
	if err != nil {
		log.Logger.Warn().Msgf("Error closing file")
		return err
	}
	return nil
}
