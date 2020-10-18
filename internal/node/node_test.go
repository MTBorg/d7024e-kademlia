package node_test

import (
	"fmt"
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpc"
	"kademlia/internal/shortlist"
	"sync"

	"kademlia/internal/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupDataHandleResponses(t *testing.T) {
	addr := address.New("127.0.0.1:1234")
	targetId := kademliaid.NewRandomKademliaID()
	n := node.Node{}
	n.Init(addr)
	k := 2

	// k total contacts in RT
	for i := 0; i < k; i++ {
		c := contact.NewContact(kademliaid.NewRandomKademliaID(), addr)
		n.RoutingTable.AddContact(c)
	}

	sl := shortlist.NewShortlist(n.ID, n.FindKClosest(n.ID, nil, k))
	channels := make([]chan string, k)
	assert.Equal(t, k, sl.Len())

	numProbed, rpcIDs := n.ProbeAlpha(sl, &channels, "", 2)
	assert.Equal(t, k, numProbed)
	assert.Equal(t, k, len(rpcIDs))

	// Should add the contacts recieved as response to the shortlist if the value
	// was not returned
	var wg sync.WaitGroup
	wg.Add(1)
	var res string
	go func(wg *sync.WaitGroup, res *string) {
		defer wg.Done()
		// wait for response
		*res = n.LookupDataHandleResponses(sl, targetId, numProbed, &channels, rpcIDs)
	}(&wg, &res)

	// mock data and write to channel to simulate response from network
	id1 := kademliaid.NewRandomKademliaID()
	id2 := kademliaid.NewRandomKademliaID()
	resp1 := fmt.Sprintf("%s!%s", id1.String(), addr.String())
	resp2 := fmt.Sprintf("%s!%s", id2.String(), addr.String())
	channels[0] <- resp1
	channels[1] <- resp2
	wg.Wait()

	assert.Equal(t, "", res)
	assert.Equal(t, k+2, sl.Len())
	assert.Nil(t, n.RPCPool.GetEntry(rpcIDs[0]))
	assert.Nil(t, n.RPCPool.GetEntry(rpcIDs[1]))

	// Should return the value if it is found in the response
	sl = shortlist.NewShortlist(n.ID, n.FindKClosest(n.ID, nil, k))
	channels = make([]chan string, k)
	numProbed, rpcIDs = n.ProbeAlpha(sl, &channels, "", 2)
	wg.Add(1)
	go func(wg *sync.WaitGroup, res *string) {
		defer wg.Done()
		*res = n.LookupDataHandleResponses(sl, targetId, numProbed, &channels, rpcIDs)
	}(&wg, &res)

	// value returned in chan 2
	valResp := fmt.Sprintf("VALUE=hello world/SENDERID=%s", id2.String())
	channels[0] <- resp1
	channels[1] <- valResp
	wg.Wait()

	assert.Equal(t, fmt.Sprintf("hello world, from node: %s", id2.String()), res)

}

func TestLookupContactHandleResponses(t *testing.T) {
	addr := address.New("127.0.0.1:1234")
	targetId := kademliaid.NewRandomKademliaID()
	n := node.Node{}
	n.Init(addr)
	k := 2

	// k total contacts in RT
	for i := 0; i < k; i++ {
		c := contact.NewContact(kademliaid.NewRandomKademliaID(), addr)
		n.RoutingTable.AddContact(c)
	}

	sl := shortlist.NewShortlist(n.ID, n.FindKClosest(n.ID, nil, k))
	channels := make([]chan string, k)
	assert.Equal(t, k, sl.Len())

	numProbed, rpcIDs := n.ProbeAlpha(sl, &channels, "", 2)
	assert.Equal(t, k, numProbed)
	assert.Equal(t, k, len(rpcIDs))

	// Should add the contacts recieved as response to the shortlist
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		// wait for response
		n.LookupContactHandleResponses(sl, targetId, numProbed, &channels, rpcIDs)
	}(&wg)

	// mock data and write to channel to simulate response from network
	id1 := kademliaid.NewRandomKademliaID()
	id2 := kademliaid.NewRandomKademliaID()
	resp1 := fmt.Sprintf("%s!%s", id1.String(), addr.String())
	resp2 := fmt.Sprintf("%s!%s", id2.String(), addr.String())
	channels[0] <- resp1
	channels[1] <- resp2
	wg.Wait()

	assert.Equal(t, k+2, sl.Len())
	assert.Nil(t, n.RPCPool.GetEntry(rpcIDs[0]))
	assert.Nil(t, n.RPCPool.GetEntry(rpcIDs[1]))
}

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
