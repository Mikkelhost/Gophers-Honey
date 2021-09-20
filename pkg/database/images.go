package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
)

type Image struct {
	GUID primitive.ObjectID `bson:"_id,omitempty"`
	Name string `bson:"name"json:"name"`
	Id uint32 `bson:"image_id"json:"image_id"`
}

func GetImages() ([]Image, error) {
	var imageList []Image

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	results, err := db.Database(DB_NAME).Collection(DB_IMAGE_COLL).Find(ctx, bson.M{})

	if err != nil {
		log.Logger.Warn().Msgf("Error retrieving image list: %s", err)
		return nil, err
	}

	for results.Next(ctx) {
		var image Image

		if err = results.Decode(&image); err != nil {
			log.Logger.Warn().Msgf("Error decoding result: %s", err)
			return nil, err
		}

		imageList = append(imageList, image)
	}

	for _, image := range imageList {
		log.Logger.Debug().Msgf("Found image with image ID: %d, ip: %s", image.Id, image.Name)
	}

	return imageList, nil
}

// NewImage puts image info into the db for the user to fetch for the frontend
func NewImage(image Image) (uint32, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	imageId := createRandID("image_id", DB_IMAGE_COLL)
	image.Id = imageId
	_, err := db.Database(DB_NAME).Collection(DB_IMAGE_COLL).InsertOne(ctx, image)

	if err != nil {
		return 0, err
	}

	return imageId, nil
}

func RemoveImage(Image) error {
	return nil
}