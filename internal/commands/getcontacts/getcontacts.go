package getcontacts

import (
	"github.com/rs/zerolog/log"
	"kademlia/internal/node"
)

type GetContacts struct{}

func (g *GetContacts) Execute() (string, error) {
	log.Debug().Msg("Executing getcontacts command")
	return node.KadNode.RoutingTable.GetContacts(), nil
}

func (g *GetContacts) ParseOptions(options []string) error {
	return nil
}

func (g *GetContacts) PrintUsage() {
	log.Info().Msg("Usage: getcontacts")
}
