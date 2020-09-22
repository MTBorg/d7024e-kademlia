package addcontact_test

import (
	// "fmt"
	"kademlia/internal/address"
	"kademlia/internal/commands/addcontact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	var addcCmd *addcontact.AddContact
	var err error

	// should not return an error if a  nodeid and address is specified
	addcCmd = new(addcontact.AddContact)
	err = addcCmd.ParseOptions([]string{"nodeid", "address"})
	assert.Nil(t, err)

	// should set node ID and Address
	addcCmd = new(addcontact.AddContact)
	err = addcCmd.ParseOptions([]string{"nodeid", "address"})
	assert.Equal(t, addcCmd.Id, "nodeid")
	assert.Equal(t, addcCmd.Address, "address")

	// should return an error if an address is specified but not node ID
	addcCmd = new(addcontact.AddContact)
	err = addcCmd.ParseOptions([]string{"address"})
	assert.NotNil(t, err)

	// should return an error if a node ID is specified but not address
	addcCmd = new(addcontact.AddContact)
	err = addcCmd.ParseOptions([]string{"nodeid"})
	assert.NotNil(t, err)
}

func TestExecute(t *testing.T) {
	var addcCmd *addcontact.AddContact

	// should add the contact
	node := node.Node{}
	node.Init(address.New("127.0.0.1:1776"))
	addcCmd = new(addcontact.AddContact)
	id := kademliaid.NewRandomKademliaID().String()
	addcCmd.ParseOptions([]string{id, "127.0.0.1:1776"})
	res, err := addcCmd.Execute(&node)
	assert.Equal(t, "Contact added: 127.0.0.1:1776", res)
	assert.Nil(t, err)
}
