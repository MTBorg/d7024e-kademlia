package forget_test

import (
	"kademlia/internal/commands/forget"
	"kademlia/internal/contact"
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"kademlia/internal/nodedata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	var err error
	var n node.Node
	value := "hello world"
	hash := kademliaid.NewKademliaID(&value)
	var f forget.Forget
	contacts := &[]contact.Contact{}

	// Should not return an error if the value exists
	f = forget.Forget{}
	n = node.Node{NodeData: nodedata.NodeData{DataStore: datastore.New()}}
	n.DataStore.Insert(value, contacts, nil, nil)
	err = f.ParseOptions([]string{hash.String()})
	assert.Nil(t, err)
	_, err = f.Execute(&n)
	assert.NoError(t, err)

	// Should return an error if the value does not exists
	f = forget.Forget{}
	n = node.Node{NodeData: nodedata.NodeData{DataStore: datastore.New()}}
	err = f.ParseOptions([]string{hash.String()})
	assert.Nil(t, err)
	_, err = f.Execute(&n)
	assert.Error(t, err)
}

func TestParseOptions(t *testing.T) {
	var err error
	value := "hello world"
	hash := kademliaid.NewKademliaID(&value)
	var f forget.Forget

	// Should return an error if the hash was not specified
	f = forget.Forget{}
	err = f.ParseOptions([]string{})
	assert.EqualError(t, err, "Missing hash")

	// Should not return an error if the hash was specified
	f = forget.Forget{}
	err = f.ParseOptions([]string{hash.String()})
	assert.NoError(t, err)
}

func TestPrintUsage(t *testing.T) {
	f := forget.Forget{}
	assert.Equal(t, "USAGE: forget <hash>", f.PrintUsage())
}
