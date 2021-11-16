package hpfeedsHandler

import (
	"errors"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/d1str0/hpfeeds"
)

type Identities struct {
	IDs []hpfeeds.Identity
}

func NewDB() *Identities {
	i := hpfeeds.Identity{
		Ident:       "test_ident",
		Secret:      "12345",
		SubChannels: []string{"opencanary_events"},
		PubChannels: []string{"opencanary_events"},
	}
	t := &Identities{IDs: []hpfeeds.Identity{i}}
	return t
}

func (t *Identities) Identify(ident string) (*hpfeeds.Identity, error) {
	if ident == "test_ident" {
		return &t.IDs[0], nil
	}
	return nil, errors.New("identifier: Unknown identity")
}

func Broker() error {
	db := NewDB()

	broker := &hpfeeds.Broker{
		Name: "broker",
		Port: 10000,
		DB:   db,
	}

	log.Logger.Info().Msgf("Starting hpfeeds broker: %s on port %d", broker.Name, broker.Port)

	err := broker.ListenAndServe()
	if err != nil {
		log.Logger.Warn().Msgf("ListenAndServe error: %s", err)
		return err
	}

	return nil
}
