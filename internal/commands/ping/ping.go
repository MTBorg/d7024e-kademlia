package ping

import (
	"errors"
	"kademlia/internal/address"
	"kademlia/internal/node"

	"github.com/rs/zerolog/log"
)

type Ping struct {
	Target string
}

func (p Ping) Execute(node *node.Node) (string, error) {
	log.Trace().Str("Target", p.Target).Msg("Executing ping command")
	adr := address.New(p.Target)
	node.Network.SendPingMessage(node.ID, adr)

	return "Ping sent!", nil

}

func (p *Ping) ParseOptions(options []string) error {
	if len(options) < 1 {
		return errors.New("Missing nodeID and target address")
	}
	p.Target = options[0]
	return nil
}

func (p *Ping) PrintUsage() string {
	return "Usage: ping {target address}"
}
