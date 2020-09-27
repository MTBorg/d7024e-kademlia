package put_test

import (
	"github.com/stretchr/testify/assert"
	"kademlia/internal/address"
	"kademlia/internal/commands/put"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"os"
	"testing"
)

func TestExecute(t *testing.T) {

	// should return an empty string
	os.Setenv("K", "20")
	node := &node.Node{}
	id := kademliaid.NewRandomKademliaID()
	adr := address.New("127.0.0.1")
	node.Init(adr)
	c := contact.NewContact(id, adr)
	node.RoutingTable.AddContact(c)

	put := put.Put{}
	err := put.ParseOptions([]string{"", "TEST"})
	assert.Nil(t, err)

	resp, err := put.Execute(node)
	assert.Nil(t, err)
	assert.Equal(t, "", resp)

}

func TestParseOptions(t *testing.T) {
	var putCmd *put.Put
	var err error

	// should not return an error if content is specified
	putCmd = new(put.Put)
	err = putCmd.ParseOptions([]string{"address", "content"})
	assert.Nil(t, err)

	// should return an error
	putCmd = new(put.Put)
	err = putCmd.ParseOptions([]string{})
	assert.NotNil(t, err)

}
