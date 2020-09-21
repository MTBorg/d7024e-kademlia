package node_test

import (
	"kademlia/internal/address"
	"kademlia/internal/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	addr := address.New("address")
	node.KadNode.Init(addr)

	// should initialize the node variables
	assert.NotNil(t, node.KadNode.RoutingTable)
}
