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

func createRandUuid() uint32 {
	rand.Seed(time.Now().Unix())
	uuid := rand.Uint32()
	for isDeviceInCollection(uuid, "uuid", DB_DEV_COLL) {
		uuid = rand.Uint32()
		log.Info().Msg("Running \"while loop\"")
	}
	return uuid
}

func ipToInt(ipStr string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ipStr).To4()), binary.BigEndian, &long)
	return long
}

func getContextWithTimeout() (context.Context, context.CancelFunc) {
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return context, cancel
}
