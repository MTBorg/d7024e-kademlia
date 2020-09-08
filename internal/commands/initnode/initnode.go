package initnode

import (
	"errors"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"kademlia/internal/routingtable"

	"github.com/rs/zerolog/log"
)

type InitNode struct {
	Address string
}

// Initialize the node by generating a NodeID and creating a new routing table
// containing itself as a contact
func (i *InitNode) Execute() (string, error) {
	log.Debug().Msg("Executing init command")
	log.Info().Msg("Initializing node...")

	id := kademliaid.NewRandomKademliaID()
	me := contact.NewContact(id, i.Address)
	node.KadNode = node.Node{
		Id:           id,
		RoutingTable: routingtable.NewRoutingTable(me),
	}

	log.Info().Str("NodeID", id.String()).Msg("ID assigned")

	return "Node initialized", nil
}

func (i *InitNode) ParseOptions(options []string) error {
	if len(options) < 1 {
		return errors.New("Missing address")
	}

	i.Address = options[0]
	return nil
}

func (i *InitNode) PrintUsage() {
	log.Info().Msg("Usage: init {address}")
}
