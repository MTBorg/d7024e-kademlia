package exit

import (
	"kademlia/internal/node"
	"os"

	"github.com/rs/zerolog/log"
)

type Exit struct {
}

// this is just so we can test execute
var ExitFunction = os.Exit

func (e Exit) Execute(node *node.Node) (string, error) {
	log.Trace().Msg("Executing exit command")
	log.Info().Msg("Node exiting...")
	ExitFunction(0)
	return "Node exited", nil
}

func (e *Exit) ParseOptions(options []string) error {
	return nil
}

func (e *Exit) PrintUsage() string {
	return "Usage: exit "
}
