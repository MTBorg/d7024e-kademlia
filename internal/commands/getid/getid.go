package getid

import (
	"kademlia/internal/node"

	"github.com/rs/zerolog/log"
)

type GetId struct {
}

// getid returns the nodes kademlia ID
func (g GetId) Execute() (string, error) {
	log.Debug().Msg("Executing getid command")
	return node.KadNode.Id.String(), nil
}

func (g *GetId) ParseOptions(options []string) error {
	return nil
}

func (g *GetId) PrintUsage() {
	log.Info().Msg("Usage: getid")
}