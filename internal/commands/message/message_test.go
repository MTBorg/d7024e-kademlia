package message_test

import (
	"kademlia/internal/address"
	"kademlia/internal/commands/message"
	"kademlia/internal/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	node := &node.Node{}
	node.Init(address.New("127.0.0.1:1234"))
	msg := message.Message{}

	// should be able to send a message to a valid target
	msg.ParseOptions([]string{"127.0.0.1:1337", "hejsan"})
	resp, err := msg.Execute(node)
	assert.Nil(t, err)
	assert.Equal(t, "Message sent!", resp)

	// should not be able to send a message to an invalid target
	msg.ParseOptions([]string{"123", "hello"})
	resp, err = msg.Execute(node)
	assert.Error(t, err)
}

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

func TestPrintUsage(t *testing.T) {
	// should be equal
	var msgCmd *message.Message
	assert.Equal(t, msgCmd.PrintUsage(), "Usage: msg {target address} {message content}")

}
