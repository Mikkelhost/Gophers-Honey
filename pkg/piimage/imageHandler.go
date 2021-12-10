package piimage

import (
	"errors"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strconv"

	diskfs "github.com/diskfs/go-diskfs"
)

//InsertConfig
//Makes a copy from a base image and inserts a config into the /boot folder of the image
func InsertConfig(conf model.PiConf, id uint32) error {
	yaml, err := yaml.Marshal(&conf)
	if err := copyImage(id); err != nil {
		log.Logger.Warn().Msgf("Error copying base image: %s", err)
		return err
	}
	disk, err := diskfs.Open("images/" + strconv.FormatUint(uint64(id), 10) + ".img")
	if err != nil {
		log.Logger.Warn().Msgf("Error opening img: %s", err)
		return err
	}
	fs, err := disk.GetFilesystem(1) // assuming it is the whole disk, so partition = 0
	if err != nil {
		log.Logger.Warn().Msgf("Error getting file system: %s", err)
		return err
	}
	file, err := fs.OpenFile("/config.yml", os.O_CREATE|os.O_RDWR)
	if err != nil {
		log.Logger.Warn().Msgf("Error creating config file in img: %s", err)
		return err
	}
	if _, err := file.Write(yaml); err != nil {

	}
	return nil
}

//copyImage
//Copies the base image into a new identical image with an id in the name.
func copyImage(id uint32) error {
	if _, err := os.Stat("images/" + strconv.FormatUint(uint64(id), 10) + ".img"); os.IsNotExist(err) {
		log.Logger.Info().Msg("Creating image")
		baseFileStat, err := os.Stat("images/base.img")
		if err != nil {
			log.Logger.Warn().Msgf("Error opening basestat image for copy: %s", err)
			return err
		}
		if !baseFileStat.Mode().IsRegular() {
			log.Logger.Warn().Msg("images/base.img is not a regular file")
			return err
		}
		base, err := os.Open("images/base.img")
		if err != nil {
			log.Logger.Warn().Msgf("Error Opening base image for copy: %s", err)
			return err
		}
		custom, err := os.Create("images/" + strconv.FormatUint(uint64(id), 10) + ".img")
		if err != nil {
			log.Logger.Warn().Msgf("Error creating custom image destination: %s", err)
			return err
		}
		defer custom.Close()
		nBytes, err := io.Copy(custom, base)
		if err != nil {
			log.Logger.Warn().Msgf("Error copying base to custom: %s", err)
			return err
		}
		log.Logger.Info().Msgf("Copied %s bytes!", nBytes)
		return nil
	} else {
		log.Logger.Info().Msgf("Image already exists")
		return errors.New("image already exists")
	}
}

// DeleteImage finds the specified image file and deletes it from disk.
func DeleteImage(imageID uint32) error {
	filePath := "images/" + strconv.FormatUint(uint64(imageID), 10) + ".img"

	if _, err := os.Stat(filePath); err == nil{
		log.Logger.Info().Msg("Deleting image")
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
		return nil
	} else {
		log.Logger.Warn().Msgf("Image not found")
		return errors.New("image not found")
	}
}


