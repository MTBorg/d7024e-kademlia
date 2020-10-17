package get_test

import (
	"kademlia/internal/commands/get"
	"kademlia/internal/contact"
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"kademlia/internal/nodedata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOption(t *testing.T) {
	var g get.Get
	var options []string
	var err error
	// should not return an error if the hash was provided
	g = *new(get.Get)
	options = []string{"somehash"}
	err = g.ParseOptions(options)
	assert.NoError(t, err)

	// should return an error if hash was not provided
	g = *new(get.Get)
	options = []string{}
	err = g.ParseOptions(options)
	assert.Error(t, err)
}

func TestExecute(t *testing.T) {
	// TODO: Not tested since .net lib

	// should return the value if it existed locally
	contacts := &[]contact.Contact{}
	value := "hello world"
	hash := kademliaid.NewKademliaID(&value)
	n := node.Node{NodeData: nodedata.NodeData{DataStore: datastore.New()}}
	n.DataStore.Insert(value, contacts, nil, nil)

	cmd := get.Get{}
	cmd.ParseOptions([]string{hash.String()})
	res, err := cmd.Execute(&n)
	assert.NoError(t, err)
	assert.Equal(t, "hello world, from local node", res)
}

func TestPrintUsage(t *testing.T) {
	// should be equal
	var getCmd *get.Get
	assert.Equal(t, getCmd.PrintUsage(), "USAGE: get <hash>")

}
