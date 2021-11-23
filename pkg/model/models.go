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
	SSH    bool `bson:"ssh" yaml:"ssh" json:"ssh"`
	FTP    bool `bson:"ftp" yaml:"ftp" json:"ftp"`
	TELNET bool `bson:"telnet" yaml:"telnet" json:"telnet"`
	HTTP   bool `bson:"http" yaml:"http" json:"http"`
	SMB    bool `bson:"smb" yaml:"smb" json:"smb"`
}

// Configuration struct matches a device ID with enabled services. Is only
// used when retrieving configuration data from the database.
type Configuration struct {
	DeviceID  uint32  `bson:"device_id,omitempty" json:"device_id"`
	NICVendor string  `bson:"nic_vendor" json:"nic_vendor"`
	Hostname  string  `bson:"hostname" json:"hostname"`
	Services  Service `bson:"services" json:"services"`
}

// Device struct is used to specify device information.
type Device struct {
	GUID       primitive.ObjectID `bson:"_id,omitempty" json:"guid"`
	DeviceID   uint32             `bson:"device_id,omitempty" json:"device_id"`
	IP         uint32             `bson:"ip,omitempty"`
	IpStr      string             `bson:"ip_str,omitempty" json:"ip_str"`
	Hostname   string             `bson:"hostname" json:"hostname"`
	Configured bool               `bson:"configured" json:"configured"`
	Services   Service            `bson:"services" json:"services"`
	LastSeen   time.Time          `bson:"last_seen" json:"last_seen"`
}

type Log struct {
	GUID         primitive.ObjectID `bson:"_id,omitempty"`
	DeviceID     uint32             `bson:"device_id,omitempty" json:"device_id"`
	LogID        uint32             `bson:"log_id,omitempty" json:"log_id"`
	DstHost      string             `bson:"dst_host" json:"dst_host"`
	DstPort      uint16             `bson:"dst_port" json:"dst_port"`
	SrcHost      string             `bson:"src_host" json:"src_host"`
	SrcPort      uint16             `bson:"src_port" json:"src_port"`
	LogTimeStamp time.Time          `bson:"log_time_stamp" json:"log_time_stamp"`
	Message      string             `bson:"message,omitempty" json:"message"`
	Level        int                `bson:"level" json:"level"`
	LogType      string             `bson:"log_type" json:"log_type"` // TODO: Possibly no use or opencanary specific.
	RawLog       string             `bson:"raw_log" json:"raw_log"`
}

// Log severity levels.
var (
	CRITICAL      = 0
	SCAN          = 1
	INFORMATIONAL = 2
)

var LogLevelMap = map[int]string{
	0: "Critical",
	1: "Scan",
	2: "Informational",
}

type Image struct {
	GUID        primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	DateCreated string             `bson:"date_created" json:"date_created"`
	Id          uint32             `bson:"image_id" json:"image_id"`
}

type DeviceAuth struct {
	DeviceId  uint32 `json:"device_id,omitempty"`
	DeviceKey string `json:"device_key,omitempty"`
}

type DBUser struct {
	FirstName            string `bson:"first_name" json:"first_name"`
	LastName             string `bson:"last_name" json:"last_name"`
	Email                string `bson:"email" json:"email"`
	Role                 string `bson:"role" json:"role"`
	NotificationsEnabled bool   `bson:"notifications_enabled" json:"notifications_enabled"`
	Username             string `bson:"username" json:"username"`
	UsernameLower        string `bson:"username_lower" json:"username_lower"`
	PasswordHash         string `bson:"password_hash,omitempty" json:"password_hash,omitempty"`
}

var (
	UserRole  = "User"
	AdminRole = "Admin"
)

/* API Call related structs.
 */

type APIResponse struct {
	Error string `json:"error"`
}

type APIUser struct {
	FirstName            string `json:"firstName,omitempty"`
	LastName             string `json:"lastName,omitempty"`
	Email                string `json:"email,omitempty"`
	Role                 string `json:"role,omitempty"`
	NotificationsEnabled bool   `json:"notifications_enabled"`
	CurrPassword         string `json:"curr_password,omitempty"`
	Username             string `json:"username,omitempty"`
	Password             string `json:"password,omitempty"`
	ConfirmPw            string `json:"confirmPw,omitempty"`
	Token                string `json:"token,omitempty"`
}

type Claims struct {
	Authorized bool   `json:"authorized"`
	Exp        int64  `json:"exp"`
	Role       string `json:"role"`
	Username   string `json:"username"`
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
	C2        string `json:"c2"`
	Port      string `json:"port"`
}

type PiConfResponse struct {
	Status    string  `json:"status"`
	DeviceId  uint32  `json:"device_id"`
	NICVendor string  `json:"nic_vendor"`
	Hostname  string  `json:"hostname"`
	Services  Service `json:"services"`
}

/* RaspberryPi image related structs.
 */

type Heartbeat struct {
	DeviceID  uint32    `json:"device_id"`
	TimeStamp time.Time `json:"time_stamp"`
}

type PiConf struct {
	C2         string  `yaml:"c2"`
	IpStr      string  `yaml:"ip_str"`
	Hostname   string  `yaml:"hostname"`
	NICVendor  string  `yaml:"nic_vendor"`
	Mac        string  `yaml:"mac"`
	Configured bool    `yaml:"configured"`
	Port       int     `yaml:"port"`
	DeviceID   uint32  `yaml:"device_id"`
	DeviceKey  string  `yaml:"device_key"`
	Services   Service `yaml:"services"`
}

/* Notification related structs.
 */

type SmtpServer struct {
	Username string `bson:"username" yaml:"username"`
	Password string `bson:"password" yaml:"password"`
	SmtpHost string `bson:"smtp_host" yaml:"smtp_host"`
	SmtpPort uint16 `bson:"smtp_port" yaml:"smtp_port"`
}

/* Service configuration related structs.
 */

type Config struct {
	Configured bool       `yaml:"configured"`
	SmtpServer SmtpServer `yaml:"smtp_server"`
}
