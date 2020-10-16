package node_test

import (
	"fmt"
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpc"
	"kademlia/internal/shortlist"

	"kademlia/internal/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeserializeContacts(t *testing.T) {
	addr := address.New("127.0.0.1:1234")
	targetId := kademliaid.NewRandomKademliaID()
	id1 := kademliaid.NewRandomKademliaID()
	id2 := kademliaid.NewRandomKademliaID()
	data := fmt.Sprintf("%s!%s %s!%s", id1.String(), addr.String(), id2.String(), addr.String())

	// should deserialize the contacts when data contains correctly formatted
	// contacts
	res := node.DeserializeContacts(data, targetId)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, *id1, *res[0].ID)
	assert.Equal(t, *id2, *res[1].ID)

	// should return an empty array of contacts if the data is empty
	res = node.DeserializeContacts("", targetId)
	assert.NotNil(t, res)
	assert.Equal(t, 0, len(res))
}

func TestProbeAlpha(t *testing.T) {
	addr := address.New("127.0.0.1:1234")
	n := node.Node{}
	n.Init(addr)
	alpha := 3
	k := 5

	// k total contacts in RT
	for i := 0; i < k; i++ {
		c := contact.NewContact(kademliaid.NewRandomKademliaID(), addr)
		n.RoutingTable.AddContact(c)
	}

	sl := shortlist.NewShortlist(n.ID, n.FindKClosest(n.ID, nil, 5))
	channels := make([]chan string, k)
	assert.Equal(t, k, sl.Len())

	// should probe the alpha closest contacts if alpha unprobed contacts exist
	numProbed, rpcIDs := n.ProbeAlpha(sl, &channels, "", alpha)
	assert.Equal(t, alpha, numProbed)
	assert.Equal(t, alpha, len(rpcIDs))
	for i := 0; i < k; i++ {
		if i < alpha {
			assert.True(t, sl.Entries[i].Probed)
		} else {
			assert.False(t, sl.Entries[i].Probed)
		}
	}
	for i := 0; i < k; i++ {
		if i < alpha {
			assert.NotNil(t, channels[i])
		} else {
			assert.Nil(t, channels[i])
		}
	}

	// should probe fewer than alpha contacts if not enough unprobed contacts
	// exist
	numProbed, rpcIDs = n.ProbeAlpha(sl, &channels, "", alpha)
	assert.Equal(t, k-alpha, numProbed)
	assert.Equal(t, k-alpha, len(rpcIDs))
	for i := 0; i < k; i++ {
		assert.True(t, sl.Entries[i].Probed)
	}

	// should not probe any contacts if no unprobed contacts exists
	numProbed, rpcIDs = n.ProbeAlpha(sl, &channels, "", alpha)
	assert.Equal(t, 0, numProbed)
	assert.Equal(t, 0, len(rpcIDs))
}

func TestStore(t *testing.T) {
	n := node.Node{}
	d := datastore.New()
	n.DataStore = d

	// should be equal
	value := "TEST"
	id := kademliaid.NewKademliaID(&value)
	contacts := &[]contact.Contact{}
	n.Store(&value, contacts, nil)
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

func TestFindKClosest(t *testing.T) {
	n := node.Node{}
	addr := address.New("127.0.1.1")
	n.Init(addr)
	key := kademliaid.FromString("ffffffffffffffffffffffffffffffffffffffff")
	id1 := kademliaid.FromString("fffffffffffffffffffffffffffffffffffffff0")
	id2 := kademliaid.FromString("ffffffffffffffffffffffffffffffffffffff00")
	id3 := kademliaid.FromString("fffffffffffffffffffffffffffffffffffff000")
	c1 := contact.NewContact(id1, addr)
	c2 := contact.NewContact(id2, addr)
	c3 := contact.NewContact(id3, addr)
	n.RoutingTable.AddContact(c1)
	n.RoutingTable.AddContact(c2)
	n.RoutingTable.AddContact(c3)

	// should return the k closest contacts to the key without returning any
	// contact with the same ID as the requestorID
	kClosest := n.FindKClosest(key, id1, 3)
	// contact c1 should not be returned since is has the same id as the requestor
	assert.Equal(t, 2, len(kClosest))
}
