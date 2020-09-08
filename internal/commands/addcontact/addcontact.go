package addcontact

import (
	"errors"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"

	"github.com/rs/zerolog/log"
)

type AddContact struct {
	Id      string
	Address string
}

func (a *AddContact) Execute() (string, error) {
	log.Debug().Msg("Executing addcontact command")
	node.KadNode.RoutingTable.AddContact(contact.NewContact(kademliaid.FromString(a.Id), a.Address))
	return "Contact added", nil
}

func (a *AddContact) ParseOptions(options []string) error {
	if len(options) < 2 {
		return errors.New("Missing contact id or address")
	}

	a.Id = options[0]
	a.Address = options[1]
	return nil
}

func (a *AddContact) PrintUsage() {
	log.Info().Msg("Usage: addcontact {nodeID} {address}")
}
