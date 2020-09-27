package node_test

import (
	"kademlia/internal/address"
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpc"

	"kademlia/internal/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	n := node.Node{}
	d := datastore.New()
	n.DataStore = d

	// should be equal
	value := "TEST"
	id := kademliaid.NewKademliaID(&value)
	n.Store(&value)
	assert.Equal(t, n.NodeData.DataStore.Get(id), "TEST")
}

func TestNewRPCWithID(t *testing.T) {

	// should be equal
	senderid := kademliaid.NewRandomKademliaID()
	rpcid := kademliaid.NewRandomKademliaID()
	adr := address.New("127.0.0.1")
	nodeId := node.NewRPCWithID(senderid, "TEST", adr, rpcid)
	assert.NotNil(t, nodeId)
	assert.Equal(t, senderid, nodeId.SenderId)
	assert.Equal(t, rpcid, nodeId.RPCId)
	assert.Equal(t, adr, nodeId.Target)
	assert.Equal(t, "TEST", nodeId.Content)
}

func TestNewRPC(t *testing.T) {

	// should be equal
	n := node.Node{}
	adr := address.New("127.0.0.1")
	rpc1 := n.NewRPC("TEST", adr)
	senderId := kademliaid.NewRandomKademliaID()
	rpc2 := rpc.New(senderId, "TEST", adr)
	assert.Equal(t, rpc1.Content, rpc2.Content)
	assert.Equal(t, rpc1.Target, rpc2.Target)

}

func TestInit(t *testing.T) {
	node := node.Node{}
	adr := address.New("address")
	node.Init(adr)

	// should initialize the node variables
	assert.NotNil(t, node.RoutingTable)
}
