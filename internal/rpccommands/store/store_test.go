package store_test

import (
	"kademlia/internal/address"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"kademlia/internal/rpccommands/store"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	var s store.Store
	var options []string
	var err error

	// Should set file content if passed
	options = []string{"this is some file content"}
	fileContent := "this is some file content"
	err = s.ParseOptions(&options)
	assert.NoError(t, err)
	node := node.Node{}
	node.Init(address.New(""))
	s.Execute(&node)
	assert.Equal(t, fileContent, node.DataStore.Get(kademliaid.NewKademliaID(&fileContent)))
}

func TestParseOptions(t *testing.T) {
	var s store.Store
	var options []string
	var err error

	// Should return an error if file content has not been specified
	options = []string{}
	err = s.ParseOptions(&options)
	assert.EqualError(t, err, "Received empty STORE RPC")

	// Should set file content if passed
	options = []string{"this", "is", "some", "file", "content"}
	err = s.ParseOptions(&options)
	assert.NoError(t, err)
	assert.Equal(t, "this is some file content", reflect.ValueOf(s).FieldByName("fileContent").String())
}
