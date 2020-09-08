package main

import (
	"kademlia/internal/command/listener"
	. "kademlia/internal/contact"
	// "kademlia/internal/message"
	// "kademlia/internal/bucket"
	"kademlia/internal/network"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	log.Info().Msg("Starting node...")
	go cmdlistener.Listen()
	// msglistener.Listen()
	network.Listen(1776)

}
