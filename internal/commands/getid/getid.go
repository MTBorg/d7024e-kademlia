package getid

import (
	"kademlia/internal/node"

	"github.com/rs/zerolog/log"
)

type GetId struct {
}

// getid returns the nodes kademlia ID
func (g GetId) Execute(node *node.Node) (string, error) {
	log.Trace().Msg("Executing getid command")
	return node.ID.String(), nil
}

func (g *GetId) ParseOptions(options []string) error {
	return nil
}

func (g *GetId) PrintUsage() string {
	return "Usage: getid"
}
