package node_test

import (
	"kademlia/internal/address"
	"kademlia/internal/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	node := node.Node{}
	adr := address.New("address")
	node.Init(adr)

	// should initialize the node variables
	assert.NotNil(t, node.RoutingTable)
}
