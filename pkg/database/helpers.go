package database

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/rs/zerolog/log"
	"math/rand"
	"net"
	"time"
)

// createRandUuid pseudo-randomly generates a number which is checked
// against the device IDs currently in the collection. Returns when no
// collision is detected.
func createRandUuid() uint32 {
	rand.Seed(time.Now().Unix())
	uuid := rand.Uint32()
	for isDeviceInCollection(uuid, "uuid", DB_DEV_COLL) {
		uuid = rand.Uint32()
		log.Debug().Msg("Running \"while loop\"") //TODO: No need to keep this?
	}
	return uuid
}

// ip2int converts an IP address from its string representation to its
// integer value.
func ip2int(ipStr string) uint32 {
	var long uint32
	err := binary.Read(bytes.NewBuffer(net.ParseIP(ipStr).To4()), binary.BigEndian, &long)
	if err != nil {
		log.Warn().Msgf("Error converting IP to int: %s", err)
		return 0
	}
	return long
}

// getContextWithTimeout is used to get a timeout context used when
// communicating with MongoDB.
func getContextWithTimeout() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx, cancel
}
