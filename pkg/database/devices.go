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
	SSH    bool `bson:"ssh" yaml:"ssh"`
	FTP    bool `bson:"ftp" yaml:"ftp"`
	TELNET bool `bson:"telnet" yaml:"telnet"`
	RDP    bool `bson:"rdp" yaml:"rdp"`
	SMB    bool `bson:"smb" yaml:"smb"`
}

// Configuration struct matches a device ID with enabled services. Is only
// used when retrieving configuration data from the database.
type Configuration struct {
	DeviceID uint32  `bson:"device_id,omitempty"`
	Services Service `bson:"services"`
}

// Device struct is used to specify device information.
type Device struct {
	GUID       primitive.ObjectID `bson:"_id,omitempty"`
	DeviceID   uint32             `bson:"device_id,omitempty"`
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
func updateConfiguration(service Service, deviceID uint32) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	if isDeviceInCollection(deviceID, "device_id", DB_CONF_COLL) {

		filter := bson.M{
			"device_id": deviceID,
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
				Uint32("device_id", deviceID).
				Msgf("Error updating device config collection: %s", err)
			return err
		}
	} else {

		config := bson.M{
			"device_id": deviceID,
			"services":  service,
		}
		_, err := db.Database(DB_NAME).Collection(DB_CONF_COLL).InsertOne(ctx, config)

		if err != nil {
			log.Logger.Warn().
				Uint32("device_id", deviceID).
				Msgf("Error adding device config to config collection %s", err)
			return err
		}
	}

	return nil
}

// ConfigureDevice configures a given device to provide specified
// services. Specifically it updates the value of "services" for the
// specific device ID in both the "device_collection" and
// "configuration_collection" collections.
func ConfigureDevice(service Service, deviceId uint32) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	filter := bson.M{
		"device_id": deviceId,
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
			Uint32("device_id", deviceId).
			Msgf("Error updating device: %s", err)
		return err
	}

	err = updateConfiguration(service, deviceId)

	if err != nil {
		return err
	}

	return nil
}

// AddDevice assigns a device ID to the device and adds it to the
// database.
func AddDevice(ipStr string) (uint32, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	deviceID := createRandDeviceID()
	ip := ip2int(ipStr)
	device := Device{
		DeviceID:   deviceID,
		IpStr:      ipStr,
		Configured: false,
		IP:         ip,
	}
	_, err := db.Database(DB_NAME).Collection(DB_DEV_COLL).InsertOne(ctx, device)

	if err != nil {
		return 0, err
	}

	return deviceID, nil
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
		log.Logger.Debug().Msgf("Found device with device ID: %d, ip: %s", device.DeviceID, device.IpStr)
	}

	return deviceList, nil
}

// RemoveDevice removes a device, with the specified device ID, from the
// database.
func RemoveDevice(devideID uint32) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	if isDeviceInCollection(devideID, "device_id", DB_CONF_COLL) {
		device := Device{
			DeviceID: devideID,
		}

		_, err := db.Database(DB_NAME).Collection(DB_DEV_COLL).DeleteOne(ctx, device)

		if err != nil {
			log.Logger.Warn().Msgf("Error removing device: %s", err)
			return err
		}
	} else {
		log.Logger.Warn().Msgf("Device ID: %d not found", devideID)
		// TODO: Perhaps we need to return an error here.
	}

	return nil
}

// GetDeviceConfiguration retrieves the configuration information stored
// for a specific device.
func GetDeviceConfiguration(deviceID uint32) (Configuration, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	filter := bson.M{
		"device_id": deviceID,
	}

	var configuration Configuration

	result := db.Database(DB_NAME).Collection(DB_CONF_COLL).FindOne(ctx, filter)

	if err := result.Decode(&configuration); err != nil {
		log.Logger.Warn().Msgf("Error decoding configuration: %s", err)
		return Configuration{}, err
	}

	log.Logger.Debug().Msgf("Found configurations for device ID %d:\n"+
		"SSH enabled: %t\n"+
		"FTP enabled: %t\n"+
		"Telnet enabled: %t\n"+
		"RDP enabled: %t\n"+
		"SMB enabled: %t",
		configuration.DeviceID, configuration.Services.SSH, configuration.Services.FTP,
		configuration.Services.TELNET, configuration.Services.RDP, configuration.Services.SMB)

	return configuration, nil
}
