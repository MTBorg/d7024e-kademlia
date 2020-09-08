package main

import (
	"kademlia/internal/command/listener"
	"kademlia/internal/message/listener"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	log.Info().Msg("Starting node...")

	go cmdlistener.Listen()
	msglistener.Listen()
}
