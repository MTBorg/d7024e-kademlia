package findnode_test

import (
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpccommands/findnode"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var findNode *findnode.FindNode
	addr := address.New("address")
	requestor := contact.NewContact(kademliaid.NewRandomKademliaID(), addr)

	// should return a FindNode object
	rpcId := kademliaid.NewRandomKademliaID()
	findNode = findnode.New(&requestor, rpcId)
	assert.IsType(t, &findnode.FindNode{}, findNode)
}

func TestParseOptions(t *testing.T) {
	var findNode *findnode.FindNode
	addr := address.New("address")
	requestor := contact.NewContact(kademliaid.NewRandomKademliaID(), addr)

	rpcId := kademliaid.NewRandomKademliaID()
	findNode = findnode.New(&requestor, rpcId)

	// should not report an error if the ID is specified
	options := []string{"id"}
	assert.NoError(t, findNode.ParseOptions(&options))

	// should report an error if the ID is not specified
	options = []string{}
	assert.Error(t, findNode.ParseOptions(&options))
}
