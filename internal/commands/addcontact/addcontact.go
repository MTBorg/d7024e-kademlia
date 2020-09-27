package addcontact

import (
	"errors"
	"fmt"
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"

	"github.com/rs/zerolog/log"
)

type AddContact struct {
	Id      string
	Address string
}

func (a *AddContact) Execute(node *node.Node) (string, error) {
	log.Debug().Msg("Executing addcontact command")
	adr := address.New(a.Address)
	node.RoutingTable.AddContact(contact.NewContact(kademliaid.FromString(a.Id), adr))
	return "Contact added: " + fmt.Sprint(adr.String()), nil
}

func (a *AddContact) ParseOptions(options []string) error {
	if len(options) < 2 {
		return errors.New("Missing contact id or address")
	}

	a.Id = options[0]
	a.Address = options[1]
	return nil
}

func (a *AddContact) PrintUsage() string {
	return "Usage: addcontact {nodeID} {address}"
}
