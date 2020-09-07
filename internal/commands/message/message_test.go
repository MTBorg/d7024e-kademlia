package message_test

import (
	"kademlia/internal/commands/message"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	var msgCmd *message.Message
	var err error

	// should not return an error if an address and message is specified
	msgCmd = new(message.Message)
	err = msgCmd.ParseOptions([]string{"address", "message"})
	assert.Nil(t, err)

	// should set the specified ip as the target
	msgCmd = new(message.Message)
	err = msgCmd.ParseOptions([]string{"address", "message"})
	assert.Equal(t, msgCmd.Target, "address")
	assert.Equal(t, msgCmd.Content, "message")

	// should return an error if an address isn't specified
	msgCmd = new(message.Message)
	err = msgCmd.ParseOptions([]string{})
	assert.NotNil(t, err)

	// should return an error if an address is specified but not content
	msgCmd = new(message.Message)
	err = msgCmd.ParseOptions([]string{"address"})
	assert.NotNil(t, err)
}
