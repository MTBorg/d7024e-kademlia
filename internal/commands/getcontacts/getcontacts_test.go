package getcontacts_test

import (
	"fmt"
	"kademlia/internal/address"
	"kademlia/internal/commands/getcontacts"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	var getcsCmd *getcontacts.GetContacts

	// should return nil
	assert.Nil(t, getcsCmd.ParseOptions([]string{"123"}))
	assert.Nil(t, getcsCmd.ParseOptions([]string{""}))

}

func TestExecute(t *testing.T) {
	addr := address.New("127.0.0.1")
	var getcsCmd *getcontacts.GetContacts
	var n node.Node

	// On a non-initialized node
	// should indicate that the node is not initialized if the routing table
	// does not exist
	n = node.Node{}
	res, err := getcsCmd.Execute(&n)
	assert.Equal(t, "The node is not initilized, it does not contain a routing table or any contacts", res)
	assert.Nil(t, err)

	// On a initialized node
	n = node.Node{}
	n.Init(addr)
	// should return the contacts in the routing table
	c := contact.NewContact(kademliaid.NewKademliaIDInRange(n.ID, 0), addr)
	n.RoutingTable.AddContact(c)
	res, err = getcsCmd.Execute(&n)
	assert.Nil(t, err)
	expected := fmt.Sprintf("\nContacts:\n\nBucket 160:\ncontact(\"%s\", \"%s\")\n\nEnd of contacts.\nTotal number of contacts: 1", c.ID, c.Address.String())
	assert.Equal(t, expected, res)
}

func TestPrintUsage(t *testing.T) {
	// should be equal
	var getcsCmd *getcontacts.GetContacts
	assert.Equal(t, getcsCmd.PrintUsage(), "Usage: getcontacts")

}
