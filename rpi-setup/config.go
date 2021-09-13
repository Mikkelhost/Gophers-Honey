package main

import (
	"fmt"
	"io/ioutil"

	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"

	"gopkg.in/yaml.v3"
)

type Record struct {
	Hostname string `yaml:"SomeRandomHostname"`
	DevideID uint32 `yaml:"DeviceID"`
}

type Services struct {
	SSH bool `yaml:"SSH"`
	FTP bool `yaml:"FTP"`
}

type Config struct {
	Record   Record   `yaml:"Settings"`
	Services Services `yaml:"Services"`
}

func createConfigFile() {

	log.Logger.Info().Msgf("[*]\tConfiguring")

	devideID := api_call_addDevice()
	config := Config{Record: Record{Hostname: "ChangedToSomething", DevideID: devideID}, Services: Services{SSH: false, FTP: false}}
	data, err := yaml.Marshal(&config)

	if err != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err)
	}
	err2 := ioutil.WriteFile("config.yaml", data, 0644)

	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err2)
	}
	log.Logger.Info().Msgf("[+]\tConfiguring [DONE]")
	readConfigFile()
}

func readConfigFile() {

	yfile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err)
	}

	data := make(map[string]Services)
	err2 := yaml.Unmarshal(yfile, &data)
	if err2 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err)
	}

	data2 := make(map[uint32]Services)
	err3 := yaml.Unmarshal(yfile, &data2)
	if err3 != nil {
		log.Logger.Error().Msgf("[X]\tError - ", err)
	}

	for k, v := range data {
		fmt.Printf("%s: %t \n", k, v)
	}
	log.Logger.Info().Msgf("[*]\tDeviceID -> %d", err3)
}
