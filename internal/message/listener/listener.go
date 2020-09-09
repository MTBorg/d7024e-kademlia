package msglistener

import (
	"kademlia/internal/contact"
	"kademlia/internal/message"
	"net"

	"github.com/rs/zerolog/log"
)

const port = ":1776"

func received(c *net.UDPConn) {
	for {
		buf := make([]byte, 512)
		nr, addr, err := c.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		data := buf[0:nr]
		msg, err := message.Deserialize(string(data))
		if err == nil {
			log.Info().Str("Content", msg.Content).Str("SenderId", msg.SenderId.String()).Msg("Received message")

			c := contact.NewContact(msg.SenderId, addr.String())
			log.Debug().Str("Id", c.ID.String()).Str("Address", c.Address).Msg("Updating bucket")
			// TODO: Add to routing table
		} else {
			log.Warn().Str("Error", err.Error()).Msg("Failed to deserialize message")
		}
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
