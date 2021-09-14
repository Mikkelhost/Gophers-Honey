package config

import (
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Configured bool `yaml:"configured,omitempty"`
}

var conf *Config
func CreateConfFile(){
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

func GetServiceConfig () (*Config, error){
	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Logger.Warn().Msgf("Error reading file: %s",err)
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Logger.Warn().Msgf("Error unmarshalling yaml: %s", err)
		return nil, err
	}
	conf = &config
	return &config, nil
}
