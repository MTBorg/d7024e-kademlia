package exit

import (
	"kademlia/internal/node"
	"os"

	"github.com/rs/zerolog/log"
)

type Exit struct {
}

func (e Exit) Execute(node *node.Node) (string, error) {
	log.Debug().Msg("Executing exit command")
	log.Info().Msg("Node exiting...")
	os.Exit(0)
	return "Node exited", nil
}

func (e *Exit) ParseOptions(options []string) error {
	return nil
}

func (e *Exit) PrintUsage() {
	log.Info().Msg("Usage: exit ")
}
