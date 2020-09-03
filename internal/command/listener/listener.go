package cmdlistener

import (
	"github.com/rs/zerolog/log"
	"kademlia/internal/command/parser"
	"net"
	"os"
)

// Clears the socket at socketAddress
//
// Useful for making sure a socket does not try to connect to an already
// bound address.
// Works by removing the file specified by the specified socket address.
func ClearSocket(socketAddress string) {
	log.Debug().Str("SocketAddress", socketAddress).Msg("Clearing socket")
	if err := os.RemoveAll(socketAddress); err != nil {
		log.Error().Str("SocketAddress", socketAddress).Msg("Failed to clear socket")
	} else {
		log.Debug().Str("SocketAddress", socketAddress).Msg("Socket cleared")
	}
}

func respond(c net.Conn) {
	buf := make([]byte, 512)
	nr, err := c.Read(buf)
	if err != nil {
		return
	}

	data := buf[0:nr]
	log.Info().Str("Command", string(data)).Msg("Received command")

	command := cmdparser.ParseCmd(string(data))

	// Execute command
	var executionResult string
	if command != nil {
		executionResult, err = command.Execute()

		// Write response
		if err == nil {
			log.Debug().Str("Message", executionResult).Msg("Sending response")
			_, err = c.Write([]byte(executionResult))
			if err != nil {
				log.Error().Msgf("Failed to write response: %s", err)
			}
		}
	}

	c.Close()
}

func Listen() {
	const socketAddress = "/tmp/echo.sock"

	ClearSocket(socketAddress)

	l, err := net.Listen("unix", socketAddress)
	if err != nil {
		log.Error().Msgf("Failed to listen: %s", err)
	}
	log.Info().Str("Address", socketAddress).Msg("Listening on socket")
	defer l.Close()

	for {
		c, err := l.Accept()
		if err == nil {
			log.Info().Str("Address", socketAddress).Msg("Received message from socket")
			go respond(c)
		} else {
			log.Error().Msgf("Listener failed to accept: %s", err)
		}
	}
}
