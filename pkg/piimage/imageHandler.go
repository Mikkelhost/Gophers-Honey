package piimage

import (
	"github.com/Mikkelhost/Gophers-Honey/pkg/database"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"gopkg.in/yaml.v2"
	"io"
	"os"

	diskfs "github.com/diskfs/go-diskfs"
)


type PiConf struct {
	HostName string `yaml:"hostname"`
	DeviceID uint32 `yaml:"device_id"`
	Services database.Service `yaml:"services"`
}

func InsertConfig(conf PiConf) {
	conf.HostName = "localhost"
	conf.DeviceID = 0
	yaml, err := yaml.Marshal(&conf)
	copyImage()
	disk, err := diskfs.Open("images/custom.img")
	if err != nil {
		log.Logger.Warn().Msgf("Error opening img: %s", err)
	}
	fs, err := disk.GetFilesystem(1) // assuming it is the whole disk, so partition = 0
	if err != nil {
		log.Logger.Warn().Msgf("Error getting file system: %s", err)
	}
	//files, err := fs.ReadDir("/") // this should list everything
	/*if err != nil {
		log.Panic(err)
	}*/
	file, err := fs.OpenFile("/config.yml", os.O_CREATE|os.O_RDWR)
	if err != nil {
		log.Logger.Warn().Msgf("Error creating config file in img: %s", err)
	}
	/*b := make([]byte, 1024, 1024)
	file.Read(b)
	fmt.Println(string(b))*/
	// if err != nil {
	//     log.Panic(err)
	// }
	file.Write(yaml)
}

func copyImage(){
	if _, err := os.Stat("images/custom.img"); os.IsNotExist(err) {
		log.Logger.Info().Msg("Creating image")
		baseFileStat, err := os.Stat("images/base.img")
		if err != nil {
			log.Logger.Fatal().Msgf("Error opening basestat image for copy: %s", err)
			return
		}
		if !baseFileStat.Mode().IsRegular() {
			log.Logger.Fatal().Msg("images/base.img is not a regular file")
			return
		}
		base, err := os.Open("images/base.img")
		if err != nil {
			log.Logger.Fatal().Msgf("Error Opening base image for copy: %s", err)
			return
		}
		custom, err := os.Create("images/custom.img")
		if err != nil {
			log.Logger.Fatal().Msgf("Error creating custom image destination: %s", err)
			return
		}
		defer custom.Close()
		nBytes, err := io.Copy(custom, base)
		if err != nil {
			log.Logger.Fatal().Msgf("Error copying base to custom: %s", err)
			return
		}
		log.Logger.Info().Msgf("Copied %s bytes!", nBytes)
	}else {
		log.Logger.Info().Msg("Image already exists")
	}
}
