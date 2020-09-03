package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"kademlia/internal/command/listener"
	. "kademlia/internal/contact"
	"kademlia/internal/message"
	"os"
	"time"
)

type Kademlia struct {
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	log.Info().Msg("Starting node...")

	go cmdlistener.Listen()
	msglistener.Listen()
}
