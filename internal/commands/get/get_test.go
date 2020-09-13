package get_test

import (
	"kademlia/internal/commands/get"
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
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
	var g get.Get
	var res string

	// should not return a value if it does not exist
	g = *new(get.Get)
	g.ParseOptions([]string{"non existent hash"})
	res, _ = g.Execute()
	assert.Equal(t, res, "Key not found")

	//should return the value if it does exist
	g = *new(get.Get)
	message := "some message"
	datastore.Store.Insert(message)
	id := kademliaid.NewKademliaID(&message)
	g.ParseOptions([]string{(&id).String()})
	res, _ = g.Execute()
	assert.Equal(t, res, "some message")
}
