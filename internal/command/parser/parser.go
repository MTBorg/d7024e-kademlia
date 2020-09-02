package cmdparser

import (
	"strings"

	"github.com/rs/zerolog/log"
	. "kademlia/internal/command"
	"kademlia/internal/commands/ping"
)

func ParseCmd(s string) Command {
	fields := strings.Fields(s)

	var command Command

	// Assume the string has already been checked to contain a command
	switch cmd := fields[0]; cmd {
	case "ping":
		command = new(ping.Ping)
	default:
		log.Error().Str("command", cmd).Msg("Received unknown command")
		return nil
	}

	err := command.ParseOptions(fields[1:]) // Extract all options
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("Failed to parse options")
		command.PrintUsage()
		return nil
	}

	return command
}
