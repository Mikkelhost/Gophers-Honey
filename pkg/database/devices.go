package database

import (
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Service struct is used to specify enabled/disabled services in a
// configuration.
type Service struct {
	SSH    bool `bson:"ssh"`
	FTP    bool `bson:"ftp"`
	TELNET bool `bson:"telnet"`
	RDP    bool `bson:"rdp"`
	SMB    bool `bson:"smb"`
}

// Device struct is used to specify device information.
type Device struct {
	GUID       primitive.ObjectID `bson:"_id,omitempty"`
	UUID       uint32             `bson:"uuid,omitempty"`
	IP         uint32             `bson:"ip,omitempty"`
	IpStr      string             `bson:"ip_str,omitempty" json:"ip_str"`
	Configured bool               `bson:"configured"`
	Services   Service            `bson:"services"`
}

// isDeviceInCollection reports whether a document with the specified
// device ID occurs in the given collection.
func isDeviceInCollection(value uint32, key, collection string) bool {
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

// updateConfiguration updates the device configuration data contained in
// the "configuration_collection" collection.
func updateConfiguration(service Service, uuid uint32) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	if isDeviceInCollection(uuid, "uuid", DB_CONF_COLL) {
		filter := bson.M{
			"uuid": uuid,
		}
		config := bson.M{
			"services": service,
		}
		update := bson.M{
			"$set": config,
		}
		_, err := db.Database(DB_NAME).Collection(DB_CONF_COLL).UpdateOne(ctx, filter, update)
		if err != nil {
			log.Logger.Warn().
				Uint32("uuid", uuid).
				Msgf("Error updating device config collection: %s", err)
		}
	} else {
		config := bson.M{
			"uuid":     uuid,
			"services": service,
		}
		_, err := db.Database(DB_NAME).Collection(DB_CONF_COLL).InsertOne(ctx, config)

		if err != nil {
			log.Logger.Warn().
				Uint32("uuid", uuid).
				Msgf("Error adding device config to config collection %s", err)
		}
	}
}

// ConfigureDevice configures a given device to provide specified
// services. Specifically it updates the value of "services" for the
// specific device ID in both the "device_collection" and
// "configuration_collection" collections.
func ConfigureDevice(service Service, uuid uint32) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()
	filter := bson.M{
		"uuid": uuid,
	}
	config := Device{
		Configured: true,
		Services:   service,
	}
	update := bson.M{
		"$set": config,
	}
	_, err := db.Database(DB_NAME).Collection(DB_DEV_COLL).UpdateOne(ctx, filter, update)

	if err != nil {
		log.Logger.Warn().
			Uint32("uuid", uuid).
			Msgf("Error updating device: %s", err)
	}
	updateConfiguration(service, uuid)
}

// AddDevice assigns a device ID to the device and adds it to the
// database.
func AddDevice(ipStr string) (uint32, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	uuid := createRandUuid()
	ip := ip2int(ipStr)
	device := Device{
		UUID:       uuid,
		IpStr:      ipStr,
		Configured: false,
		IP:         ip,
	}
	_, err := db.Database(DB_NAME).Collection(DB_DEV_COLL).InsertOne(ctx, device)
	if err != nil {
		return 0, err
	}
	return uuid, nil
}

// GetAllDevices retrieves and returns a list of all devices currently in
// the database. Specifically it retrieves all devices contained in the
// "device_collection" collection.
func GetAllDevices() ([]Device, error) {
	var deviceList []Device

	ctx, cancel := getContextWithTimeout()
	defer cancel()
	results, err := db.Database(DB_NAME).Collection(DB_DEV_COLL).Find(ctx, bson.M{})
	if err != nil {
		log.Logger.Warn().Msgf("Error retrieving device list: %s", err)
		return nil, err
	}

	for results.Next(ctx) {
		var device Device
		if err = results.Decode(&device); err != nil {
			log.Logger.Warn().Msgf("Error decoding result: %s", err)
			return nil, err
		}
		deviceList = append(deviceList, device)
	}
	for _, device := range deviceList {
			log.Logger.Debug().Msgf("Found device with uuid: %d, ip: %s", device.UUID, device.IpStr)
	}
	return deviceList, nil
}

// RemoveDevice removes the given device from collection,
// if the given uuid is valid
func RemoveDevice(uuid uint32) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	if isDeviceInCollection(uuid, "uuid", DB_CONF_COLL) {
		device := Device{
			UUID: uuid,
		}
		_, err := db.Database(DB_NAME).Collection(DB_DEV_COLL).DeleteOne(ctx, device)
		if err != nil {
			log.Logger.Warn().Msgf("Error removing device: %s", err)
			return
		}
	} else {
		log.Logger.Warn().Msgf("Device ID: %d not found", uuid)
	}

}
