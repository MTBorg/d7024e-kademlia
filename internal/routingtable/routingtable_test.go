package routingtable_test

import (
	"kademlia/internal/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContacts(t *testing.T) {
	// should return message informing that the routingtable is empty
	assert.Equal(t, "Empty! Please, populate the routingtable...", node.KadNode.RoutingTable.GetContacts())

}
