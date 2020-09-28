package get

import (
	"errors"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"

	"github.com/rs/zerolog/log"
)

type Get struct {
	hash kademliaid.KademliaID
}

func (get *Get) Execute(node *node.Node) (string, error) {
	log.Trace().Msg("Executing get command")

	// Check local storage
	value := node.DataStore.Get(get.hash)
	if value == "" {
		log.Debug().Str("Key", get.hash.String()).Msg("Value not found locally")

		value = node.LookupData(&get.hash)
	}

	if value == "" {
		return "", errors.New("Key not found")
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

func (get *Get) PrintUsage() string {
	return "USAGE: get <hash>"
}
