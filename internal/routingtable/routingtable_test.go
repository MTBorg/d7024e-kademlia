package routingtable_test

import (
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"kademlia/internal/routingtable"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoutingTable(t *testing.T) {

	// should not be nil and return empty string of contacts
	id := kademliaid.NewRandomKademliaID()
	adr := address.New("127.0.0.1")
	c := contact.NewContact(id, adr)
	rt := routingtable.NewRoutingTable(c)
	assert.NotNil(t, rt)
}

func TestAddContact(t *testing.T) {
	// node := node.Node{}
	id := kademliaid.NewRandomKademliaID()
	adr := address.New("127.0.0.1")
	c := contact.NewContact(id, adr)
	rt := routingtable.NewRoutingTable(c)

	// should be empty string
	rt.AddContact(c)
	assert.Equal(t, "\nContacts:\n\nEnd of contacts.\nTotal number of contacts: 0", rt.GetContacts())

	// should not be empty string
	id = kademliaid.NewRandomKademliaID()
	c = contact.NewContact(id, adr)
	rt.AddContact(c)
	assert.NotEqual(t, "", rt.GetContacts())

}

func TestFindClosestContacts(t *testing.T) {
	node := node.Node{}
	id := kademliaid.NewRandomKademliaID()
	id2 := kademliaid.NewRandomKademliaID()
	adr := address.New("127.0.0.1")
	c := contact.NewContact(id, adr)
	node.Init(adr)
	node.RoutingTable.AddContact(c)

	// should not return a index
	assert.NotNil(t, node.RoutingTable.FindClosestContacts(id, id2, 1))

}

func TestGetContacts(t *testing.T) {
	node := node.Node{}
	// should return message informing that the routingtable does not exist
	assert.Equal(t, "The node is not initilized, it does not contain a routing table or any contacts", node.RoutingTable.GetContacts())

	// should return a contact in newline string format
	id := kademliaid.NewRandomKademliaID()
	adr := address.New("127.0.0.1")
	c := contact.NewContact(id, adr)
	node.Init(adr)
	node.RoutingTable.AddContact(c)
	contacts := node.RoutingTable.GetContacts()
	assert.NotEqual(t, "Empty! Please, populate the routingtable...", contacts)

}

func TestgetBucketIndex(t *testing.T) {
	node := node.Node{}
	id := kademliaid.NewRandomKademliaID()

	// should not return a index
	assert.Nil(t, node.RoutingTable.GetBucketIndex(id))

	// should return a index
	adr := address.New("127.0.0.1")
	c := contact.NewContact(id, adr)
	node.Init(adr)
	node.RoutingTable.AddContact(c)
	assert.NotNil(t, node.RoutingTable.GetBucketIndex(id))
}
