package getcontacts_test

import (
	"kademlia/internal/commands/getcontacts"
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
	var getcsCmd *getcontacts.GetContacts

	// should indicate that the node is not initialized if the routing table
	// does not exist
	node := node.Node{}
	res, err := getcsCmd.Execute(&node)
	assert.Equal(t, "The node is not initilized, it does not contain a routing table or any contacts", res)
	assert.Nil(t, err)
}

func TestPrintUsage(t *testing.T) {
	// should be equal
	var getcsCmd *getcontacts.GetContacts
	assert.Equal(t, getcsCmd.PrintUsage(), "Usage: getcontacts")

}
