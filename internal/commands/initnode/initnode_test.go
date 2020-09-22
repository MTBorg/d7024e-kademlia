package initnode_test

import (
	"testing"

	"kademlia/internal/commands/initnode"
	"kademlia/internal/node"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	var initCmd *initnode.InitNode
	var err error

	// should not return an error if an address is specified
	initCmd = new(initnode.InitNode)
	err = initCmd.ParseOptions([]string{"address"})
	assert.Nil(t, err)

	// should set the specified ip as the address
	initCmd = new(initnode.InitNode)
	err = initCmd.ParseOptions([]string{"address"})
	assert.Equal(t, initCmd.Address, "address")

	// should return an error if an address isn't specified
	initCmd = new(initnode.InitNode)
	err = initCmd.ParseOptions([]string{})
	assert.NotNil(t, err)
}

func TestExecute(t *testing.T) {
	var initCmd *initnode.InitNode

	// should initialize the node
	initCmd = new(initnode.InitNode)
	initCmd.ParseOptions([]string{"address"})
	node := node.Node{}
	res, err := initCmd.Execute(&node)
	assert.Equal(t, "Node initialized", res)
	assert.Nil(t, err)
}
