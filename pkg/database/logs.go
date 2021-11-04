package database

import (
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	ttlIndexName = "time_stamp_1"
	ttlIndexSet  = false
)

// AddLog assigns a log with timestamp and message tied to a device ID and adds it to the
// database. Also sets a time to live index (if not set) of 3 months. This ensures that
// logs are deleted from the database after 3 months.
// TODO: timestamp needs proper implementation + message
func AddLog(deviceID uint32, timeStamp time.Time, message string) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	if !ttlIndexSet {
		indexNames, err := getIndexNames(DB_LOG_COLL)
		if err != nil {
			return err
		}
		if stringAppearsInArray(ttlIndexName, indexNames) {
			ttlIndexSet = true
		} else {
			err = setTTLIndex(7889231) // 7889231 seconds = 3 months default expiration date
			if err != nil {
				return err
			}
			ttlIndexSet = true
		}
	}

	logID := createRandID("log_id", DB_LOG_COLL)
	deviceLog := model.Log{
		DeviceID:  deviceID,
		LogID:     logID,
		TimeStamp: timeStamp,
		Message:   message,
	}
	_, err := db.Database(DB_NAME).Collection(DB_LOG_COLL).InsertOne(ctx, deviceLog)

	if err != nil {
		return err
	}

	return nil
}

// GetAllLogs retrieves and returns a list of all logs currently in
// the database.
func GetAllLogs() ([]model.Log, error) {
	var logList []model.Log

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	results, err := db.Database(DB_NAME).Collection(DB_LOG_COLL).Find(ctx, bson.M{})

	if err != nil {
		log.Logger.Warn().Msgf("Error retrieving log list: %s", err)
		return nil, err
	}

	for results.Next(ctx) {
		var deviceLog model.Log

		if err = results.Decode(&deviceLog); err != nil {
			log.Logger.Warn().Msgf("Error decoding result: %s", err)
			return nil, err
		}

		logList = append(logList, deviceLog)
	}

	for _, deviceLog := range logList {
		log.Logger.Debug().Msgf("Found log with log ID: %d", deviceLog.LogID)
	}

	return logList, nil
}

// GetLog gets a single log from the database based on the given logID
func GetLog(logID uint32) (model.Log, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	filter := bson.M{
		"logID": logID,
	}
	var deviceLog model.Log

	result := db.Database(DB_NAME).Collection(DB_LOG_COLL).FindOne(ctx, filter)

	if err := result.Decode(&deviceLog); err != nil {
		log.Logger.Warn().Msgf("Error decoding result: %s", err)
		return model.Log{}, err
	}
	return deviceLog, nil
}

// RemoveLog removes a log, with the specified ID, from the
// database.
func RemoveLog(logID uint32) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	if isIdInCollection(logID, "log_id", DB_LOG_COLL) {
		deviceLog := model.Log{
			LogID: logID,
		}

		_, err := db.Database(DB_NAME).Collection(DB_LOG_COLL).DeleteOne(ctx, deviceLog)

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

// setTTLIndex sets a TTL or "expireAfter" index on the "time_stamp" field
// of the logs in the collection. This allows mongoDB to automatically
// remove logs, which are older than the specified TTL, from the
// collection.
func setTTLIndex(seconds int32) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	ttlIndex := mongo.IndexModel{
		Keys: bson.M{
			"time_stamp": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(seconds),
	}

	_, err := db.Database(DB_NAME).Collection(DB_LOG_COLL).Indexes().CreateOne(ctx, ttlIndex)
	if err != nil {
		log.Logger.Warn().Msgf("Error creating TTL index: %s", err)
		return err
	}

	return nil
}

// UpdateTTLIndex updates the "setExpireAfterSeconds" index of the
// "log_collection" collection by removing and resetting the index.
func UpdateTTLIndex(seconds int32) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	_, err := db.Database(DB_NAME).Collection(DB_LOG_COLL).Indexes().DropOne(ctx, ttlIndexName)

	if err != nil {
		log.Logger.Warn().Msgf("Error removing TTL index: &s", err)
		return err
	}

	err = setTTLIndex(seconds)

	if err != nil {
		return err
	}

	return nil
}
