package ping_test

import (
	"kademlia/internal/address"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpccommands/ping"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	adr := address.New("127.0.0.1:1776")
	p := ping.New(adr, kademliaid.NewRandomKademliaID())
	options := []string{"hello", "abc"}
	//Should never return an error
	assert.NoError(t, p.ParseOptions(&options))
}
