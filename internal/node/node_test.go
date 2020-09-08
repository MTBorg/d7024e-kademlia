package node_test

import (
	"kademlia/internal/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	node.KadNode.Init("address")

	// should initialize the node variables
	assert.NotNil(t, node.KadNode.Id)
	assert.NotNil(t, node.KadNode.RoutingTable)
}
