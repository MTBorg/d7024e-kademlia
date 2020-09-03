package msglistener

import (
	"net"

	"github.com/rs/zerolog/log"
)

const port = ":1776"

func received(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			continue
		}

		data := buf[0:nr]
		log.Info().Str("Content", string(data)).Msg("Received message with content,")

	}
}

// Listen for UDP on port
func Listen() {
	addr, err := net.ResolveUDPAddr("udp4", port)
	if err != nil {
		log.Error().Msgf("Failed to resolve UDP Address: %s", err)
	}

	ln, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Error().Msgf("Failed to listen on UDP Address: %s", err)
	}
	log.Info().Str("Address", addr.String()).Msg("Listening on UDP packets on address")
	defer ln.Close()

	received(ln)

}
