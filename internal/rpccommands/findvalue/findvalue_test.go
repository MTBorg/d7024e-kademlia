package findvalue_test

import (
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpccommands/findvalue"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	addr := address.New("127.0.0.1")
	reqID := kademliaid.NewRandomKademliaID()
	req := contact.NewContact(reqID, addr)
	rpcID := kademliaid.NewRandomKademliaID()

	// should return a new FIND_VALUE RPC
	find := findvalue.New(&req, rpcID)
	assert.IsType(t, findvalue.FindValue{}, *find)
}

func TestParseOptions(t *testing.T) {
	var err error
	find := findvalue.FindValue{}

	//Should return an error if id was not specified
	err = find.ParseOptions(&[]string{})
	assert.EqualError(t, err, "Missing hash")

	//Should not return an error if id was not specified
	err = find.ParseOptions(&[]string{"someid"})
	assert.NoError(t, err)
}
