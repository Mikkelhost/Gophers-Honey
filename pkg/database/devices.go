package database

import (
	"errors"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/Mikkelhost/Gophers-Honey/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// setDefaultConfiguration sets a default configuration when a new RPi is connected.
func setDefaultConfiguration(deviceID uint32) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	services := model.Service{}
	configuration := model.Configuration{
		DeviceID: deviceID,
		Services: services,
	}

	_, err := db.Database(DB_NAME).Collection(DB_CONF_COLL).InsertOne(ctx, configuration)
	if err != nil {
		log.Logger.Warn().Msgf("Error adding default configuration to config collection")
		return err
	}
	return nil
}

// updateConfiguration updates the device configuration data contained in
// the "configuration_collection" collection.
func updateConfiguration(config model.Configuration) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	if isIdInCollection(config.DeviceID, "device_id", DB_CONF_COLL) {

		filter := bson.M{
			"device_id": config.DeviceID,
		}
		conf := bson.M{
			"nic_vendor": config.NICVendor,
			"hostname":   config.Hostname,
			"services":   config.Services,
		}
		update := bson.M{
			"$set": conf,
		}

		_, err := db.Database(DB_NAME).Collection(DB_CONF_COLL).UpdateOne(ctx, filter, update)

		if err != nil {
			log.Logger.Warn().
				Uint32("device_id", config.DeviceID).
				Msgf("Error updating device config collection: %s", err)
			return err
		}
	} else {

		conf := bson.M{
			"device_id": config.DeviceID,
			"services":  config.Services,
		}
		_, err := db.Database(DB_NAME).Collection(DB_CONF_COLL).InsertOne(ctx, conf)

		if err != nil {
			log.Logger.Warn().
				Uint32("device_id", config.DeviceID).
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
func ConfigureDevice(config model.Configuration) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	filter := bson.M{
		"device_id": config.DeviceID,
	}
	update := bson.M{
		"$set": bson.M{
			"configured": true,
			"hostname":   config.Hostname,
			"services":   config.Services,
			"nic_vendor": config.NICVendor,
		},
	}
	if isIdInCollection(config.DeviceID, "device_id", DB_DEV_COLL) {
		_, err := db.Database(DB_NAME).Collection(DB_DEV_COLL).UpdateOne(ctx, filter, update)

		if err != nil {
			log.Logger.Warn().
				Uint32("device_id", config.DeviceID).
				Msgf("Error updating device: %s", err)
			return err
		}

		err = updateConfiguration(config)

		if err != nil {
			return err
		}
	} else {
		log.Logger.Warn().Msgf("Device ID: %d not found", config.DeviceID)
		return errors.New("device ID not found")
	}

	return nil
}

// AddDevice assigns a device ID to the device and adds it to the
// database.
func AddDevice(ipStr string) (uint32, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	deviceID := createRandID("device_id", DB_DEV_COLL)
	ip := ip2int(ipStr)
	device := model.Device{
		DeviceID:   deviceID,
		IpStr:      ipStr,
		Configured: false,
		IP:         ip,
	}
	_, err := db.Database(DB_NAME).Collection(DB_DEV_COLL).InsertOne(ctx, device)

	if err != nil {
		return 0, err
	}

	err = setDefaultConfiguration(deviceID)
	if err != nil {
		return 0, err
	}

	return deviceID, nil
}

// GetAllDevices retrieves and returns a list of all devices currently in
// the database. Specifically it retrieves all devices contained in the
// "device_collection" collection.
func GetAllDevices() ([]model.Device, error) {
	var deviceList []model.Device

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	results, err := db.Database(DB_NAME).Collection(DB_DEV_COLL).Find(ctx, bson.M{})

	if err != nil {
		log.Logger.Warn().Msgf("Error retrieving device list: %s", err)
		return nil, err
	}

	for results.Next(ctx) {
		var device model.Device

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
func RemoveDevice(deviceID uint32) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	if isIdInCollection(deviceID, "device_id", DB_DEV_COLL) {
		device := bson.M{
			"device_id": deviceID,
		}

		_, err := db.Database(DB_NAME).Collection(DB_DEV_COLL).DeleteOne(ctx, device)

		if err != nil {
			log.Logger.Warn().Msgf("Error removing device: %s", err)
			return err
		}

		_, err = db.Database(DB_NAME).Collection(DB_CONF_COLL).DeleteOne(ctx, device)
		if err != nil {
			log.Logger.Warn().Msgf("Error removing device: %s", err)
			return err
		}
	} else {
		log.Logger.Warn().Msgf("Device ID: %d not found", deviceID)
		return errors.New("device ID not found")
	}

	return nil
}

// GetDeviceConfiguration retrieves the configuration information stored
// for a specific device.
func GetDeviceConfiguration(deviceID uint32) (model.Configuration, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	filter := bson.M{
		"device_id": deviceID,
	}

	var configuration model.Configuration

	result := db.Database(DB_NAME).Collection(DB_CONF_COLL).FindOne(ctx, filter)

	if err := result.Decode(&configuration); err != nil {
		log.Logger.Warn().Msgf("Error decoding configuration: %s", err)
		return model.Configuration{}, err
	}

	log.Logger.Debug().Msgf("Found configurations for device ID %d:\n"+
		"Hostname: %s\n"+
		"NIC vendor: %s\n"+
		"SSH enabled: %t\n"+
		"FTP enabled: %t\n"+
		"Telnet enabled: %t\n"+
		"HTTP enabled: %t\n"+
		"SMB enabled: %t",
		configuration.DeviceID, configuration.Hostname, configuration.NICVendor, configuration.Services.SSH, configuration.Services.FTP,
		configuration.Services.TELNET, configuration.Services.HTTP, configuration.Services.SMB)

	return configuration, nil
}

// HandleHeartbeat retrieves a timestamp from the API and sets/updates the
// "last_seen" field for a given device.
func HandleHeartbeat(heartbeat model.Heartbeat) error {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	timestamp := time.Now()

	filter := bson.M{
		"device_id": heartbeat.DeviceID,
	}

	config := bson.M{
		"last_seen": timestamp,
		"ip_str":    heartbeat.IpStr,
		"ip":        ip2int(heartbeat.IpStr),
	}

	update := bson.M{
		"$set": config,
	}

	if isIdInCollection(heartbeat.DeviceID, "device_id", DB_DEV_COLL) {
		_, err := db.Database(DB_NAME).Collection(DB_DEV_COLL).UpdateOne(ctx, filter, update)

		if err != nil {
			log.Logger.Warn().
				Uint32("device_id", heartbeat.DeviceID).
				Msgf("Error updating device: %s", err)
			return err
		}
	} else {
		log.Logger.Warn().Msgf("Device ID: %d not found", heartbeat.DeviceID)
		return errors.New("device ID not found")
	}
	log.Logger.Debug().Uint32("device_id", heartbeat.DeviceID).Msgf("Successfully added timestamp: %v to device", timestamp)
	return nil
}
