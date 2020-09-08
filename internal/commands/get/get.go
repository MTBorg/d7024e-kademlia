package get

import (
	"errors"
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"

	"github.com/rs/zerolog/log"
)

type Get struct {
	hash kademliaid.KademliaID
}

func (get *Get) Execute() (string, error) {
	log.Debug().Msg("Executing get command")

	// Check local storage
	value := datastore.Store.Get(get.hash)
	if value == "" {
		log.Debug().Str("Key", get.hash.String()).Msg("Value not found locally")

		// TODO: Send FIND_NODE RPC

		return "Key not found", nil
	}

	return value, nil
}

func (get *Get) ParseOptions(options []string) error {
	if len(options) < 1 {
		return errors.New("Missing hash")
	}
	get.hash = *kademliaid.FromString(options[0])
	return nil
}

func (get *Get) PrintUsage() {
	log.Info().Msg("USAGE: get <hash>")
}