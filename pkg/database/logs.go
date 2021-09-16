package database

import (
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

var (
	ttlIndexName = "time_stamp_1"
	ttlIndexSet  = false
)

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
			err = setTTLIndex(7889231) // 7889231 seconds = 3 month expiration date
			if err != nil {
				return err
			}
			ttlIndexSet = true
		}
	}

	logID := createRandLogID()
	deviceLog := Log{
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
		var deviceLog Log

		if err = results.Decode(&deviceLog); err != nil {
			log.Logger.Warn().Msgf("Error decoding result: %s", err)
			return nil, err
		}

		logList = append(logList, deviceLog)
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
	var deviceLog Log

	result := db.Database(DB_NAME).Collection(DB_LOG_COLL).FindOne(ctx, filter)

	if err := result.Decode(&deviceLog); err != nil {
		log.Logger.Warn().Msgf("Error decoding result: %s", err)
		return Log{}, err
	}
	return deviceLog, nil
}

// RemoveLog removes a log, with the specified ID, from the
// database.
func RemoveLog(logID uint32) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	if isLogInCollection(logID, "log_id", DB_LOG_COLL) {
		deviceLog := Log{
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
		return err
	}

	return nil
}
