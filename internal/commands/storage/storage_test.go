package storage_test

import (
	"fmt"
	"kademlia/internal/address"
	"kademlia/internal/commands/storage"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	var storageCmd *storage.Storage
	var n node.Node

	storageCmd = &storage.Storage{}
	n = node.Node{}
	n.Init(address.New("127.0.0.1"))
	storedVal := "this is a test"
	storedKey := kademliaid.NewKademliaID(&storedVal)
	contacts := &[]contact.Contact{}
	n.DataStore.Insert(storedVal, contacts, nil, nil)
	expected := fmt.Sprintf("map( \n %x=%s\n)", storedKey, storedVal)

	// should return the nodes stored key-value pairs
	res, err := storageCmd.Execute(&n)
	assert.Nil(t, err)
	assert.Equal(t, expected, res)
}

func TestParseOptions(t *testing.T) {
	var storageCmd *storage.Storage

	// should not require any arguments
	storageCmd = &storage.Storage{}
	err := storageCmd.ParseOptions([]string{})
	assert.Nil(t, err)

	// should ignore any specified arguments
	storageCmd = &storage.Storage{}
	err = storageCmd.ParseOptions([]string{"hej hej hej"})
	assert.Nil(t, err)
}

func TestPrintUsage(t *testing.T) {
	var storageCmd *storage.Storage

	// should return a string
	storageCmd = &storage.Storage{}
	res := storageCmd.PrintUsage()
	assert.Equal(t, "Usage: storage", res)
}
