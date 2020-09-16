package ping_test

import (
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpccommands/ping"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	address := "someaddress"
	p := ping.New(&address, kademliaid.NewRandomKademliaID())
	options := []string{"hello", "abc"}
	//Should never return an error
	assert.NoError(t, p.ParseOptions(&options))
}
