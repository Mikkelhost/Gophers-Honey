package hpfeedsHandler

import (
	"flag"
	log "github.com/Mikkelhost/Gophers-Honey/pkg/logger"
	"github.com/d1str0/hpfeeds"
	"time"
)

func Subscribe(id, ch, secret string) error {
	var (
		host    string
		port    int
		ident   string
		auth    string
		channel string
	)
	flag.StringVar(&host, "host", "127.0.0.1", "target host")
	flag.IntVar(&port, "port", 10000, "hpfeeds port")
	flag.StringVar(&ident, "ident", id, "ident username")
	flag.StringVar(&auth, "secret", secret, "ident secret")
	flag.StringVar(&channel, "channel", ch, "channel to subscribe to")
	flag.Parse()

	// Create hpfeeds client.
	hp := hpfeeds.NewClient(host, port, ident, auth)
	hp.Log = false

	// Make go routine to extract messages.
	msgs := make(chan hpfeeds.Message)
	go func() {
		for msg := range msgs {
			log.Logger.Info().Msgf("Received message: %s, from: %s", string(msg.Payload), msg.Name)
		}
	}()

	for {
		log.Logger.Debug().Msgf("Attempting to connect to hpfeeds broker: %s on port %d using ident %s "+
			"with auth %s", host, port, ident, auth)

		// Connect to broker.
		err := hp.Connect()
		if err != nil {
			log.Logger.Warn().Msgf("Error connecting to broker server.")
			return err
		}

		log.Logger.Debug().Msgf("Successfully connected.")

		// Subscribe to channel. NB! Channel name needs to exist on broker! Else a SIGSEGV will occur!
		hp.Subscribe(channel, msgs)

		log.Logger.Debug().Msgf("Subscribing to channel: %s", channel)

		// Wait for disconnect
		<-hp.Disconnected
		log.Logger.Info().Msgf("Disconnected, attempting reconnect in 10 seconds...")
		time.Sleep(10 * time.Second)
	}
}
