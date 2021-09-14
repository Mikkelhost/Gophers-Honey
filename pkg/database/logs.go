package database

import (
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Log struct {
	GUID      primitive.ObjectID `bson:"_id,omitempty"`
	DeviceID  uint32             `bson:"device_id,omitempty" json:"device_id"`
	LogID     uint32             `bson:"log_id,omitempty" json:"log_id"`
	TimeStamp time.Time          `bson:"time_stamp,omitempty" json:"time_stamp"`
	Message   string             `bson:"message,omitempty" json:"message"`
}

// isDeviceInCollection reports whether a document with the specified
// log ID occurs in the given collection.
func isLogInCollection(value uint32, key, collection string) bool {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	filter := bson.M{
		key: value,
	}

	countOptions := options.Count().SetLimit(1)
	count, err := db.Database(DB_NAME).Collection(collection).CountDocuments(ctx, filter, countOptions)

	if err != nil {
		log.Logger.Warn().Msgf("Error counting documents: %s", err)
	}

	if count > 0 {
		return true
	}

	return false
}

// AddLog assigns a log with timestamp and message tied to a device ID and adds it to the
// database.
// TODO: timestamp needs proper implementation + message
func AddLog(deviceID uint32, timeStamp time.Time, message string) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	logID := createRandLogID()
	dlog := Log{
		DeviceID:  deviceID,
		LogID:     logID,
		TimeStamp: timeStamp,
		Message:   message,
	}
	_, err := db.Database(DB_NAME).Collection(DB_LOG_COLL).InsertOne(ctx, dlog)

	if err != nil {
		return err
	}

	return nil
}

// GetAllLogs retrieves and returns a list of all logs currently in
// the database.
func GetAllLogs() ([]Log, error) {
	var logList []Log

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	results, err := db.Database(DB_NAME).Collection(DB_LOG_COLL).Find(ctx, bson.M{})

	if err != nil {
		log.Logger.Warn().Msgf("Error retrieving log list: %s", err)
		return nil, err
	}

	for results.Next(ctx) {
		var dlog Log

		if err = results.Decode(&dlog); err != nil {
			log.Logger.Warn().Msgf("Error decoding result: %s", err)
			return nil, err
		}

		logList = append(logList, dlog)
	}

	for _, dlog := range logList {
		log.Logger.Debug().Msgf("Found log with log ID: %d", dlog.LogID)
	}

	return logList, nil
}

// GetLog gets a single log from the database based on the given logID
func GetLog(logID uint32) (Log, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	filter := bson.M{
		"logID": logID,
	}
	var dlog Log

	result := db.Database(DB_NAME).Collection(DB_LOG_COLL).FindOne(ctx, filter)

	if err := result.Decode(&dlog); err != nil {
		log.Logger.Warn().Msgf("Error decoding result: %s", err)
		return Log{}, err
	}
	return dlog, nil
}

// RemoveLog removes a log, with the specified ID, from the
// database.
func RemoveLog(logID uint32) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	if isLogInCollection(logID, "log_id", DB_LOG_COLL) {
		dlog := Log{
			LogID: logID,
		}

		_, err := db.Database(DB_NAME).Collection(DB_LOG_COLL).DeleteOne(ctx, dlog)

		if err != nil {
			log.Logger.Warn().Msgf("Error removing log: %s", err)
			return err
		}
	} else {
		log.Logger.Warn().Msgf("Log ID: %d not found", logID)
		// TODO: Perhaps we need to return an error here.
	}

	return nil
}
