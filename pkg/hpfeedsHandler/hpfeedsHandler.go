package hpfeedsHandler

import (
	"errors"
	"flag"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/d1str0/hpfeeds"
	"time"
)

type TestIdentity struct {
	IDs []hpfeeds.Identity
}

func NewDB() *TestIdentity {
	i := hpfeeds.Identity{
		Ident:       "test_ident",
		Secret:      "test_secret",
		SubChannels: []string{"test_channel"},
		PubChannels: []string{"test_channel"},
	}
	t := &TestIdentity{IDs: []hpfeeds.Identity{i}}
	return t
}

func (t *TestIdentity) Identify(ident string) (*hpfeeds.Identity, error) {
	if ident == "test_ident" {
		return &t.IDs[0], nil
	}
	return nil, errors.New("identifier: Unknown identity")
}

func Broker() error {
	db := NewDB()
	broker := &hpfeeds.Broker{
		Name: "test_broker",
		Port: 10000,
		DB:   db,
	}
	// broker.SetDebugLogger(log.Logger.Print)
	// broker.SetInfoLogger(log.Logger.Print)
	// broker.SetErrorLogger(log.Logger.Print)
	err := broker.ListenAndServe()
	if err != nil {
		log.Logger.Warn().Msgf("Error: %s", err)
		return err
	}
	return nil
}

func TestSubscriber() error {
	var (
		host    string
		port    int
		ident   string
		auth    string
		channel string
	)
	flag.StringVar(&host, "host", "127.0.0.1", "target host")
	flag.IntVar(&port, "port", 10000, "hpfeeds port")
	flag.StringVar(&ident, "ident", "test_ident", "ident username")
	flag.StringVar(&auth, "secret", "test_secret", "ident secret")
	flag.StringVar(&channel, "channel", "test_channel", "channel to subscribe to")
	flag.Parse()

	hp := hpfeeds.NewClient(host, port, ident, auth)
	hp.Log = false

	msgs := make(chan hpfeeds.Message)
	go func() {
		for foo := range msgs {
			log.Logger.Info().Msgf("Received message: %s, from: %s", string(foo.Payload), foo.Name)
		}
	}()

	for {
		log.Logger.Info().Msgf("Connecting to hpfeeds server.")
		err := hp.Connect()
		if err != nil {
			log.Logger.Warn().Msgf("Error connecting to broker server.")
			return err
		}

		// Subscribe to "flotest" and print everything coming in on it
		hp.Subscribe(channel, msgs)

		// Wait for disconnect
		<-hp.Disconnected
		log.Logger.Info().Msgf("Disconnected, attempting reconnect in 10 seconds...")
		time.Sleep(10 * time.Second)
	}
}
