package ping_test

import (
	"kademlia/internal/commands/ping"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
