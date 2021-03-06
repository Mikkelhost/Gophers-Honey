package database

import (
	"errors"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)


//GetImages
//Returns a list of images from the database.
func GetImages() ([]model.Image, error) {
	var imageList []model.Image

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	results, err := db.Database(DB_NAME).Collection(DB_IMAGE_COLL).Find(ctx, bson.M{})

	if err != nil {
		log.Logger.Warn().Msgf("Error retrieving image list: %s", err)
		return nil, err
	}

	for results.Next(ctx) {
		var image model.Image

		if err = results.Decode(&image); err != nil {
			log.Logger.Warn().Msgf("Error decoding result: %s", err)
			return nil, err
		}

		imageList = append(imageList, image)
	}

	for _, image := range imageList {
		log.Logger.Debug().Msgf("Found image with image ID: %d, name: %s, date_created: %s", image.Id, image.Name, image.DateCreated)
	}

	return imageList, nil
}

// NewImage puts image info into the db for the user to fetch for the frontend
func NewImage(image model.Image) (uint32, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()
	image.DateCreated = time.Now().Format(time.UnixDate)
	imageId := createRandID("image_id", DB_IMAGE_COLL)
	image.Id = imageId
	_, err := db.Database(DB_NAME).Collection(DB_IMAGE_COLL).InsertOne(ctx, image)

	if err != nil {
		return 0, err
	}

	return imageId, nil
}

// RemoveImage removes a specified image from the "image_collection"
// collection.
func RemoveImage(imageID uint32) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	if isIdInCollection(imageID, "image_id", DB_IMAGE_COLL) {
		image := bson.M{
			"image_id": imageID,
		}

		_, err := db.Database(DB_NAME).Collection(DB_IMAGE_COLL).DeleteOne(ctx, image)

		if err != nil {
			log.Logger.Warn().Msgf("Error removing image: %s", err)
			return err
		}
	} else {
		log.Logger.Warn().Msgf("Image ID: %d not found", imageID)
		return errors.New("image ID not found")
	}

	return nil
}
