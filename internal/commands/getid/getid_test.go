package getid_test

import (
	"kademlia/internal/commands/getid"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"kademlia/internal/nodedata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	var getidCmd *getid.GetId
	var err error

	// should not return an error (since no arugments are needed)
	getidCmd = new(getid.GetId)
	err = getidCmd.ParseOptions([]string{""})
	assert.Nil(t, err)
}

func TestExecute(t *testing.T) {
	var getidCmd *getid.GetId

	// should return the nodes ID
	id := kademliaid.NewRandomKademliaID()
	node := node.Node{NodeData: nodedata.NodeData{ID: id}}
	getidCmd = new(getid.GetId)
	res, err := getidCmd.Execute(&node)
	assert.Nil(t, err)
	assert.Equal(t, id.String(), res)
}

func TestPrintUsage(t *testing.T) {
	// should be equal
	var getidCmd *getid.GetId
	assert.Equal(t, getidCmd.PrintUsage(), "Usage: getid")

}
