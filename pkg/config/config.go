package config

import (
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Configured bool `yaml:"configured"`
}

var Conf *Config

func CreateConfFile() {
	// Checking if file already exists
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		f, err := os.Create("config.yml")
		if err != nil {
			log.Logger.Fatal().Msgf("Error creating config.yml: %s", err)
		}
		var config = Config{
			Configured: false,
		}
		yaml, err := yaml.Marshal(&config)
		if err != nil {
			log.Logger.Fatal().Msgf("Error creating Yaml string: %s", err)
		}
		f.Write(yaml)
		f.Close()
	}
}

func GetServiceConfig() (*Config, error) {
	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Logger.Warn().Msgf("Error reading file: %s", err)
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Logger.Warn().Msgf("Error unmarshalling yaml: %s", err)
		return nil, err
	}
	Conf = &config
	return &config, nil
}

func SetConfig(config Config) error {
	log.Logger.Debug().Msgf("Config to be set: %v", config)
	Conf = &config
	log.Logger.Debug().Msgf("conf: %v", *Conf)
	f, err := os.OpenFile("config.yml", os.O_RDWR, 0644)
	if err != nil {
		log.Logger.Warn().Msgf("Error opening config.yml: %s", err)
		return err
	}
	yaml, err := yaml.Marshal(Conf)
	if err != nil {
		log.Logger.Warn().Msgf("Error creating yaml string: %s", err)
		return err
	}
	f.Truncate(0)
	f.Seek(0, 0)
	if _, err := f.Write(yaml); err != nil {
		log.Logger.Warn().Msgf("Error writing conf", err)
		return err
	}

	f.Close()
	return nil
}

func SetSmtpServer(server model.SmtpServer) {

}
