package cmdparser

import (
	"strings"

	. "kademlia/internal/command"
	"kademlia/internal/commands/addcontact"
	"kademlia/internal/commands/exit"
	"kademlia/internal/commands/get"
	"kademlia/internal/commands/getid"
	"kademlia/internal/commands/initnode"
	"kademlia/internal/commands/message"
	"kademlia/internal/commands/ping"
	"kademlia/internal/commands/storage"

	"github.com/rs/zerolog/log"
)

func ParseCmd(s string) Command {
	fields := strings.Fields(s)

	var command Command

	// Assume the string has already been checked to contain a command
	switch cmd := fields[0]; cmd {
	case "ping":
		command = new(ping.Ping)
	case "msg":
		command = new(message.Message)
	case "exit":
		command = new(exit.Exit)
	case "storage":
		command = new(storage.Storage)
	case "get":
		command = new(get.Get)
	case "getid":
		command = new(getid.GetId)
	case "addcontact":
		command = new(addcontact.AddContact)
	case "init":
		command = new(initnode.InitNode)
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
