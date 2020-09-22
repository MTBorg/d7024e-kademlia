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

	// should return message informing that the routingtable is empty
	node := node.Node{}
	res, err := getcsCmd.Execute(&node)
	assert.Equal(t, "Empty! Please, populate the routingtable...", res)
	assert.Nil(t, err)
}
