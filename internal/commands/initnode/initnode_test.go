package initnode_test

import (
	"kademlia/internal/commands/initnode"
	"kademlia/internal/node"

	"github.com/stretchr/testify/assert"
	"testing"
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
	addr := "127.0.1.2:1776"

	// should initialize the node
	initCmd = new(initnode.InitNode)
	initCmd.ParseOptions([]string{addr})
	node := node.Node{}
	res, err := initCmd.Execute(&node)
	assert.Nil(t, err)
	assert.Equal(t, "Node initialized", res)
	// if the node was initialized the routing table should have the same address
	// as the argument of the command
	assert.Equal(t, addr, node.RoutingTable.GetMe().Address.String())
}

func TestPrintUsage(t *testing.T) {
	// should be equal
	var initCmd *initnode.InitNode
	assert.Equal(t, initCmd.PrintUsage(), "Usage: init {address}")

}
