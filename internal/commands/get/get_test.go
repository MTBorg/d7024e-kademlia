package get_test

import (
	"kademlia/internal/commands/get"
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
}

func TestPrintUsage(t *testing.T) {
	// should be equal
	var getCmd *get.Get
	assert.Equal(t, getCmd.PrintUsage(), "USAGE: get <hash>")

}
