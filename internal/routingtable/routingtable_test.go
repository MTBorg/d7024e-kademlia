package routingtable_test

import (
	"kademlia/internal/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContacts(t *testing.T) {
	node := node.Node{}
	// should return message informing that the routingtable is empty
	assert.Equal(t, "Empty! Please, populate the routingtable...", node.RoutingTable.GetContacts())

}
