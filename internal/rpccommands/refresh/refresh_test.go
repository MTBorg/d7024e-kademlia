package refresh_test

import (
	"kademlia/internal/address"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"kademlia/internal/rpccommands/refresh"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	r := refresh.Refresh{}
	n := node.Node{}
	addr := address.New("127.0.0.1:1234")
	n.Init(addr)
	val := "hello world"
	hash := kademliaid.NewKademliaID(&val)

	// should not fail if the value is not stored in the nodes datastore
	r.ParseOptions(&[]string{hash.String()})
	r.Execute(&n)
}

func TestParseOption(t *testing.T) {
	r := refresh.Refresh{}
	val := "hello world"
	hash := kademliaid.NewKademliaID(&val)

	// should parse the hash
	err := r.ParseOptions(&[]string{hash.String()})
	assert.Nil(t, err)

	// should return an error when the hash is missing
	err = r.ParseOptions(&[]string{})
	assert.NotNil(t, err)
}
