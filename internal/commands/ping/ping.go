package ping

import (
	"errors"
	"github.com/rs/zerolog/log"
	"kademlia/internal/contact"
	"kademlia/internal/network"
)

type Ping struct {
	Target string
}

func (p Ping) Execute() (string, error) {
	log.Debug().Str("Target", p.Target).Msg("Executing ping command")
	var contact = new(contact.Contact)
	contact.Address = p.Target
	network.Net.SendPingMessage(contact)
	result := "PING SENT!"

	return result, nil
}

func (p *Ping) ParseOptions(options []string) error {
	if len(options) < 1 {
		return errors.New("Missing target address")
	}
	p.Target = options[0]
	return nil
}

func (p *Ping) PrintUsage() {
	log.Info().Msg("Usage: ping {target address}")
}
