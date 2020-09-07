package exit

import (
	"github.com/rs/zerolog/log"
	"os"
)

type Exit struct {
}

func (e Exit) Execute() (string, error) {
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
