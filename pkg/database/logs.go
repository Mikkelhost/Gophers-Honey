package database

import (
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Log struct {
	GUID      primitive.ObjectID `bson:"_id, omitempty"`
	LogID     uint32             `bson:"log_id, omitempty"`
	TimeStamp time.Time          `bson:"time_stamp, omitempty"`
	Message   string             `bson:"message, omitempty"`
}

// AddLog assigns a log with timestamp and message tied to a device ID and adds it to the
// database.
// TODO: timestamp needs proper implementation
func AddLog(logID uint32, timeStamp time.Time, message string) (uint32, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	dlog := Log{
		LogID:     logID,
		TimeStamp: timeStamp,
		Message:   message,
	}
	_, err := db.Database(DB_NAME).Collection(DB_LOG_COLL).InsertOne(ctx, dlog)

	if err != nil {
		return 0, err
	}

	return logID, nil
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

/*
//TODO: work in progress
func GetLog(logID uint32) (Log, error) {
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

	return dlog, nil
}
*/

// RemoveLog removes a log, with the specified ID, from the
// database.
func RemoveLog(logID uint32) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	if isDeviceInCollection(logID, "logID", DB_CONF_COLL) {
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
