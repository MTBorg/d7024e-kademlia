package ping_test

import (
	"kademlia/internal/address"
	"kademlia/internal/commands/ping"
	"kademlia/internal/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	// should return a string and not an error
	addr := address.New("127.0.0.1:1234")
	node := &node.Node{}
	node.Init(addr)
	ping := ping.Ping{}
	resp, err := ping.Execute(node)
	assert.Nil(t, err)
	assert.Equal(t, "Ping sent!", resp)
}

func TestParseOptions(t *testing.T) {
	var pingCmd *ping.Ping
	var err error

	// should not return an error if an address is specified
	pingCmd = new(ping.Ping)
	err = pingCmd.ParseOptions([]string{"SomeIP"})
	assert.Nil(t, err)

	// should set the specified ip as the target
	pingCmd = new(ping.Ping)
	pingCmd.ParseOptions([]string{"SomeIP"})
	assert.Equal(t, pingCmd.Target, "SomeIP")

	// should return an error if an address isn't specified
	pingCmd = new(ping.Ping)
	err = pingCmd.ParseOptions([]string{})
	assert.NotNil(t, err)
}

func TestPrintUsage(t *testing.T) {
	// should be equal
	var pingCmd *ping.Ping
	assert.Equal(t, pingCmd.PrintUsage(), "Usage: ping {target address}")

}
