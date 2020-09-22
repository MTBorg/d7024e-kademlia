package ping

import (
	"errors"
	"kademlia/internal/address"
	"kademlia/internal/network"
	"kademlia/internal/node"

	"github.com/rs/zerolog/log"
)

type Ping struct {
	Target string
}

func (p Ping) Execute(node *node.Node) (string, error) {
	log.Debug().Str("Target", p.Target).Msg("Executing ping command")
	adr := address.New(p.Target)
	network.Net.SendPingMessage(node.ID, &adr)

	return "Ping sent!", nil

}

func (p *Ping) ParseOptions(options []string) error {
	if len(options) < 1 {
		return errors.New("Missing nodeID and target address")
	}
	p.Target = options[0]
	return nil
}

func (p *Ping) PrintUsage() {
	log.Info().Msg("Usage: ping {target address}")
}
