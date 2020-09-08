package storage

import (
	"kademlia/internal/datastore"

	"github.com/rs/zerolog/log"
)

type Storage struct{}

func (d Storage) Execute() (string, error) {
	log.Debug().Msg("Executing storage command")

	result := datastore.Store.EntriesAsString()

	return result, nil
}

func (p *Storage) ParseOptions(options []string) error {
	return nil
}

func (p *Storage) PrintUsage() {
	log.Info().Msg("Usage: data")
}
