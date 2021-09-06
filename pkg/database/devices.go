package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type Service struct {
	SSH    bool `bson:"ssh"`
	FTP    bool `bson:"ftp"`
	TELNET bool `bson:"telnet"`
	RDP    bool `bson:"rdp"`
	SMB    bool `bson:"smb"`
}

type Device struct {
	GUID     primitive.ObjectID `bson:"_id,omitempty"`
	UUID     uint32             `bson:"uuid,omitempty"`
	IP       uint32             `bson:"ip,omitempty"`
	IpStr    string             `bson:"ip_str,omitempty"`
	Services Service            `bson:"services,omitempty"`
}