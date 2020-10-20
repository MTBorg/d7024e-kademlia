package node_test

import (
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"kademlia/internal/rpc"
	"kademlia/internal/shortlist"

	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

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

func TestGetEnvIntVariable(t *testing.T) {
	// should return the value of the env var if it is set
	os.Setenv("TEST_VAR", "1337")
	assert.Equal(t, 1337, node.GetEnvIntVariable("TEST_VAR", 123))

	// should return the default value if the env var is not set
	assert.Equal(t, 123, node.GetEnvIntVariable("TEST_VAR2", 123))
}

func TestSetupLookupAlgorithm(t *testing.T) {
	n := node.Node{}
	addr := address.New("127.0.0.1:1234")
	n.Init(addr)
	c := contact.NewContact(kademliaid.NewRandomKademliaID(), addr)
	n.RoutingTable.AddContact(c)

	// should return the vars needed in the lookup algo
	alpha, k, sl, channels := node.SetupLookUpAlgorithm(&n, n.ID)
	assert.Equal(t, 3, alpha)
	assert.Equal(t, 5, k)
	assert.Equal(t, 1, sl.Len())
	assert.Equal(t, k, len(channels))
}

// fixture that returns a new node and shortlist as well as the contacts
// in the shortlist
func lookupHelper(k int, alpha int, targetId *kademliaid.KademliaID) (*node.Node, *shortlist.Shortlist, []contact.Contact) {
	addr := address.New("127.0.0.1:1234")
	n := node.Node{}
	n.Init(addr)

	// the original shortlist will contain k=5 contacts which are all far away
	// from the targetID
	for i := 0; i < k; i++ {
		id := kademliaid.FromString(strings.Repeat("f", 20+i) + strings.Repeat("0", 20-i))
		c := contact.NewContact(id, addr)
		n.RoutingTable.AddContact(c)
	}
	sl := shortlist.NewShortlist(targetId, n.FindKClosest(targetId, nil, alpha))

	return &n, sl, sl.GetContacts()
}

func TestLookupContact(t *testing.T) {
	addr := address.New("127.0.0.1:1234")
	targetId := kademliaid.FromString("ffffffffffffffffffffffffffffffffffffffff")

	// mock return values from SetupLookUpAlgorithm
	alpha := 3
	k := 5
	channels := make([]chan string, k)

	n, sl, slStateAtStart := lookupHelper(k, alpha, targetId)

	// mock of SetupLookUpAlgorithm
	node.SetupLookUpAlgorithm = func(n *node.Node, id *kademliaid.KademliaID) (int, int, *shortlist.Shortlist, []chan string) {
		return alpha, k, sl, channels
	}

	// perform the lookup in go routine
	var wg sync.WaitGroup
	wg.Add(1)
	var res []contact.Contact
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		// wait for response
		res = n.LookupContact(targetId)
	}(&wg)

	time.Sleep(time.Millisecond * 100)

	// DESCRIBE: Termination due to closest not updated after an iteration

	// first iteration
	// the 3 probed contacts respond with one contact each, all are closer than
	// the original contacts in the shortlist, this should result in 2 of the
	// original contacts being dropped from the shortlist
	id1 := kademliaid.FromString("fffffffffffffffffffffffffffffffffff00000")
	id2 := kademliaid.FromString("ffffffffffffffffffffffffffffffffffff0000")
	id3 := kademliaid.FromString("fffffffffffffffffffffffffffffffffffff000")
	resp1 := fmt.Sprintf("%s!%s", id1.String(), addr.String())
	resp2 := fmt.Sprintf("%s!%s", id2.String(), addr.String())
	resp3 := fmt.Sprintf("%s!%s", id3.String(), addr.String())
	channels[0] <- resp1
	channels[1] <- resp2
	channels[2] <- resp3

	time.Sleep(time.Millisecond * 100)

	// second iteration
	// only 1 of the 3 new contacts in the shortlist respond with a contact,
	// this contact is further away and will not be added to the shortlist
	// since the lookup did not find a closer contact it will terminate
	id4 := kademliaid.FromString("f000000000000000000000000000000000000000")
	resp1 = fmt.Sprintf("%s!%s", id4.String(), addr.String())
	resp2 = ""
	resp3 = ""
	channels[0] <- resp1
	channels[1] <- resp2
	channels[2] <- resp3

	time.Sleep(time.Millisecond * 100)

	// should return the k closest contacts (to the targetID) found during the
	// lookup
	wg.Wait()
	assert.NotNil(t, res)
	expected := []*kademliaid.KademliaID{}
	// the 3 new contacts and the two closest form the original shortlist should
	// be the k=5 closest found during the lookup
	expected = append(expected, id3, id2, id1, slStateAtStart[0].ID, slStateAtStart[1].ID)
	for i := 0; i < k; i++ {
		assert.Equal(t, *expected[i], *res[i].ID)
	}
}

func TestLookupData(t *testing.T) {
	addr := address.New("127.0.0.1:1234")
	targetId := kademliaid.FromString("ffffffffffffffffffffffffffffffffffffffff")
	alpha := 3
	k := 5
	channels := make([]chan string, k)

	n, sl, _ := lookupHelper(k, alpha, targetId)

	// mock of SetupLookUpAlgorithm
	node.SetupLookUpAlgorithm = func(n *node.Node, id *kademliaid.KademliaID) (int, int, *shortlist.Shortlist, []chan string) {
		return alpha, k, sl, channels
	}

	// perform the lookup in go routine
	var wg sync.WaitGroup
	wg.Add(1)
	var res string
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		// wait for response
		res = n.LookupData(targetId)
	}(&wg)

	time.Sleep(time.Millisecond * 100)

	// DESCRIBE: Termination due to response containing value

	// first iteration
	// one of the alpha=3 responses contains the value
	id1 := kademliaid.FromString("fffffffffffffffffffffffffffffffffff00000")
	id2 := kademliaid.FromString("ffffffffffffffffffffffffffffffffffff0000")
	id3 := kademliaid.FromString("fffffffffffffffffffffffffffffffffffff000")
	resp1 := fmt.Sprintf("%s!%s", id1.String(), addr.String())
	resp2 := fmt.Sprintf("VALUE=hello world/SENDERID=%s", id2.String())
	resp3 := fmt.Sprintf("%s!%s", id3.String(), addr.String())
	channels[0] <- resp1
	channels[1] <- resp2
	channels[2] <- resp3
	wg.Wait()

	// should return the value
	assert.Equal(t, fmt.Sprintf("hello world, from node: %s", id2.String()), res)

	// DESCRIBE: Termination due to closest not improving during an iteration
	n, sl, _ = lookupHelper(k, alpha, targetId)
	channels = make([]chan string, k)

	// mock of SetupLookUpAlgorithm
	node.SetupLookUpAlgorithm = func(n *node.Node, id *kademliaid.KademliaID) (int, int, *shortlist.Shortlist, []chan string) {
		return alpha, k, sl, channels
	}

	wg.Add(1)
	res = ""
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		// wait for response
		res = n.LookupData(targetId)
	}(&wg)
	time.Sleep(time.Millisecond * 100)

	// first iteration
	// the 3 probed contacts respond with one contact each, all are closer than
	// the original contacts in the shortlist, this should result in 2 of the
	// original contacts being dropped from the shortlist
	resp1 = fmt.Sprintf("%s!%s", id1.String(), addr.String())
	resp2 = fmt.Sprintf("%s!%s", id2.String(), addr.String())
	resp3 = fmt.Sprintf("%s!%s", id3.String(), addr.String())
	channels[0] <- resp1
	channels[1] <- resp2
	channels[2] <- resp3
	time.Sleep(time.Millisecond * 100)

	// second iteration
	// only 1 of the 3 new contacts in the shortlist respond with a contact,
	// this contact is further away and will not improve the closest contact
	// found but will still be added to the shortlist
	id4 := kademliaid.FromString("fffffffffffffffffffffffffffffffff0000000")
	resp1 = fmt.Sprintf("%s!%s", id4.String(), addr.String())
	resp2 = ""
	resp3 = ""
	channels[0] <- resp1
	channels[1] <- resp2
	channels[2] <- resp3
	time.Sleep(time.Millisecond * 100)

	// final round before returning the result (since closest didn't improve)
	// the only unprobed node in the shortlist responds with the value
	resp1 = fmt.Sprintf("VALUE=hello world/SENDERID=%s", id4.String())
	channels[0] <- resp1
	wg.Wait()

	// should return the value
	assert.Equal(t, fmt.Sprintf("hello world, from node: %s", id4.String()), res)

	// DESCRIBE: The value is not found
	//var slStateAtStart []contact.Contact
	//n, sl, slStateAtStart = lookupContactHelper(3, 3, targetId)
	n, sl, _ = lookupHelper(3, 3, targetId)
	channels = make([]chan string, k)

	// mock of SetupLookUpAlgorithm
	node.SetupLookUpAlgorithm = func(n *node.Node, id *kademliaid.KademliaID) (int, int, *shortlist.Shortlist, []chan string) {
		return alpha, k, sl, channels
	}

	wg.Add(1)
	res = ""
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		// wait for response
		res = n.LookupData(targetId)
	}(&wg)
	time.Sleep(time.Millisecond * 100)

	// First iteration
	// the probed contacts return no contacts/value
	channels[0] <- ""
	channels[1] <- ""
	channels[2] <- ""
	wg.Wait()

	// should return the k-closest contacts found
	assert.Contains(t, res, "Value not found")
}
