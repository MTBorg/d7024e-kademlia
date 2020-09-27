package storage

import (
	"kademlia/internal/node"

	"github.com/rs/zerolog/log"
)

type Storage struct{}

func (d Storage) Execute(node *node.Node) (string, error) {
	log.Debug().Msg("Executing storage command")

	result := node.DataStore.EntriesAsString()

	return result, nil
}

func (p *Storage) ParseOptions(options []string) error {
	return nil
}

func (p *Storage) PrintUsage() string {
	return "Usage: data"
}
