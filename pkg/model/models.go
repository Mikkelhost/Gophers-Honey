package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

/* Database related structs.
*/

// Service struct is used to specify enabled/disabled services in a
// configuration.
type Service struct {
	SSH    bool `bson:"ssh" yaml:"ssh"json:"ssh"`
	FTP    bool `bson:"ftp" yaml:"ftp"json:"ftp"`
	TELNET bool `bson:"telnet" yaml:"telnet"json:"telnet"`
	RDP    bool `bson:"rdp" yaml:"rdp"json:"rdp"`
	SMB    bool `bson:"smb" yaml:"smb"json:"smb"`
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
	DeviceID   uint32             `bson:"device_id,omitempty" json:"device_id"`
	IP         uint32             `bson:"ip,omitempty"`
	IpStr      string             `bson:"ip_str,omitempty" json:"ip_str"`
	Configured bool               `bson:"configured"`
	Services   Service            `bson:"services"`
}

type Log struct {
	GUID      primitive.ObjectID `bson:"_id,omitempty"`
	DeviceID  uint32             `bson:"device_id,omitempty" json:"device_id"`
	LogID     uint32             `bson:"log_id,omitempty" json:"log_id"`
	TimeStamp time.Time          `bson:"time_stamp,omitempty" json:"time_stamp"`
	Message   string             `bson:"message,omitempty" json:"message"`
}

type Image struct {
	GUID        primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"json:"name"`
	DateCreated string             `bson:"date_created"json:"date_created"`
	Id          uint32             `bson:"image_id"json:"image_id"`
}

type DeviceAuth struct {
	DeviceId  uint32 `json:"device_id,omitempty"`
	DeviceKey string `json:"device_key,omitempty"`
}

type DBUser struct {
	FirstName     string `bson:"first_name"json:"first_name"`
	LastName      string `bson:"last_name"json:"last_name"`
	Email         string `bson:"email"json:"email"`
	Username      string `bson:"username"json:"username"`
	UsernameLower string `bson:"username_lower"json:"username_lower"`
	PasswordHash  string `bson:"password_hash,omitempty"json:"password_hash,omitempty"`
}

/* API Call related structs.
 */

type APIUser struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	ConfirmPw string `json:"confirmPw,omitempty"`
	Token     string `json:"token,omitempty"`
	Error     string `json:"error"`
}

type ConfigResponse struct {
	Configured bool `json:"configured"`
}

type SetupParams struct {
	Image ImageInfo `json:"imageInfo"`
	User  APIUser   `json:"userInfo"`
}

type ImageInfo struct {
	ImageName string `json:"name"`
	Hostname  string `json:"hostname"`
	Port      string `json:"port"`
}

type PiConfResponse struct {
	Status string `json:"status"`
	DeviceId uint32 `json:"device_id"`
	Services Service `json:"services"`
}

/* RaspberryPi image related structs.
*/

type PiConf struct {
	HostName  string  `yaml:"hostname"`
	Port      int     `yaml:"port"`
	DeviceID  uint32  `yaml:"device_id"`
	DeviceKey string  `yaml:"device_key"`
	Services  Service `yaml:"services"`
}

