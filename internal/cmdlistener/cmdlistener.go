package cmdlistener

import (
	"github.com/rs/zerolog/log"
	"net"
)

func respond(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		log.Info().Str("Command", string(data)).Msg("Received command")
		_, err = c.Write(data)
		if err != nil {
			log.Error().Msgf("Failed to write response: %s", err)
		}
		c.Close()
	}
}

func Listen() {
	const socketAddress = "/tmp/echo.sock"
	l, err := net.Listen("unix", socketAddress)
	if err != nil {
		log.Error().Msgf("Failed to listen: %s", err)
	}
	log.Info().Str("Address", socketAddress).Msg("Listening on socket")
	defer l.Close()

	for {
		c, err := l.Accept()
		log.Info().Str("Address", c.LocalAddr().String()).Msg("Received message from socket")
		if err != nil {
			log.Error().Msgf("Listener failed to accept: %s", err)
		}

		go respond(c)
	}
}
