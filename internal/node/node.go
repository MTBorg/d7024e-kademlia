package node

import (
	"kademlia/internal/contact"
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
	"kademlia/internal/routingtable"

	"github.com/rs/zerolog/log"
)

type Node struct {
	Id           *kademliaid.KademliaID
	RoutingTable *routingtable.RoutingTable
}

var KadNode Node

// Initialize the node by generating a NodeID and creating a new routing table
// containing itself as a contact
func (node *Node) Init(address string) {
	id := kademliaid.NewRandomKademliaID()
	KadNode = Node{
		Id:           id,
		RoutingTable: routingtable.NewRoutingTable(contact.NewContact(id, address)),
	}
}

func (node *Node) LookupContact(target *contact.Contact) {
	// TODO
}

func (node *Node) LookupData(hash string) {
	// TODO
}

func (node *Node) Store(value *string) {
	log.Debug().Str("Value", *value).Msg("Storing value")
	datastore.Store.Insert(*value)
}
