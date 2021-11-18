package hpfeedsHandler

import (
	"errors"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/d1str0/hpfeeds"
)

type Identities struct {
	IDs []hpfeeds.Identity
}

func (t *Identities) Identify(ident string) (*hpfeeds.Identity, error) {
	for i, identity := range t.IDs {
		if ident == identity.Ident {
			return &t.IDs[i], nil
		}
	}
	return nil, errors.New("identifier: Unknown identity")
}

func NewDB() *Identities {
	i := hpfeeds.Identity{
		Ident:       "test_ident_1",
		Secret:      "12345",
		SubChannels: []string{"opencanary_events"},
		PubChannels: []string{"opencanary_events"},
	}

	j := hpfeeds.Identity{
		Ident:       "test_ident_2",
		Secret:      "54321",
		SubChannels: []string{"opencanary_events"},
		PubChannels: []string{"opencanary_events"},
	}

	backendParser := hpfeeds.Identity{
		Ident:       "backend_parser",
		Secret:      "112233",
		SubChannels: []string{"opencanary_events"},
		PubChannels: []string{"opencanary_events"},
	}

	ids := Identities{
		IDs: []hpfeeds.Identity{i, j, backendParser},
	}

	return &ids
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
