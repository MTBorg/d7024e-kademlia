package routingtable_test

import (
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"kademlia/internal/routingtable"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoutingTable(t *testing.T) {
	// should return a new routing table
	id := kademliaid.NewRandomKademliaID()
	adr := address.New("127.0.0.1")
	c := contact.NewContact(id, adr)
	rt := routingtable.NewRoutingTable(c)
	assert.NotNil(t, rt)
	assert.Equal(t, c.ID, rt.GetMe().ID)
	assert.Equal(t, 160, len(rt.Buckets))
}

func TestAddContact(t *testing.T) {
	id := kademliaid.FromString("ffffffffffffffffffffffffffffffffffffffff")
	adr := address.New("127.0.0.1")
	self := contact.NewContact(id, adr)
	rt := routingtable.NewRoutingTable(self)

	// should not be able to add the node itself as a contact
	rt.AddContact(self)
	assert.Equal(t, "\nContacts:\n\nEnd of contacts.\nTotal number of contacts: 0", rt.GetContacts())

	// should be able to add a new contact
	id = kademliaid.FromString("0000000000000000000000000000000000000000")
	c := contact.NewContact(id, adr)
	rt.AddContact(c)
	assert.Equal(t, 1, rt.GetBucket(0).Len())
}

func TestFindClosestContacts(t *testing.T) {
	adr := address.New("127.0.0.1")
	target := kademliaid.FromString("fffffffffffffffffffffffffffffffffffffff0")
	targetFurther := kademliaid.FromString("ff00000000000000000000000000000000000000")
	id1 := kademliaid.FromString("ffffffffffffffffffffffffffffffffffffff00")
	id2 := kademliaid.FromString("fffffffffffffffffffffffffffffffffffff000")
	id3 := kademliaid.FromString("ffffffffffffffffffffffffffffffffffff0000")
	id4 := kademliaid.FromString("fffffffffffffffffffffffffffffffffff00000")

	me := contact.NewContact(target, adr)
	c1 := contact.NewContact(id1, adr)
	c2 := contact.NewContact(id2, adr)
	c3 := contact.NewContact(id3, adr)
	c4 := contact.NewContact(id4, adr)

	rt := routingtable.NewRoutingTable(me)
	rt.AddContact(c1)
	rt.AddContact(c2)
	rt.AddContact(c3)
	rt.AddContact(c4)

	// Describe: Target of FindClosestContacts close to nodeID
	// it should return the k closest contacts
	kClosest := rt.FindClosestContacts(target, nil, 3)
	assert.Equal(t, 3, len(kClosest))

	// it should return the k closest contacts except any contact with the same
	// ID as the specified requestorID
	kClosest = rt.FindClosestContacts(target, id2, 5)
	assert.Equal(t, 3, len(kClosest))

	// Describe: Target of FindClosestContacts differing a lot from nodeID
	// it should return the k closest contacts
	kClosest = rt.FindClosestContacts(targetFurther, nil, 5)
	assert.Equal(t, 4, len(kClosest))

	// it should return the k closest contacts except any contact with the same
	// ID as the specified requestorID
	kClosest = rt.FindClosestContacts(targetFurther, id2, 5)
	assert.Equal(t, 3, len(kClosest))
}

func TestGetContacts(t *testing.T) {
	node := node.Node{}
	// should return message informing that the routingtable does not exist
	assert.Equal(t, "The node is not initilized, it does not contain a routing table or any contacts", node.RoutingTable.GetContacts())

	// should return a contact in newline string format
	id := kademliaid.NewRandomKademliaID()
	adr := address.New("127.0.0.1")
	c := contact.NewContact(id, adr)
	node.Init(adr)
	node.RoutingTable.AddContact(c)
	contacts := node.RoutingTable.GetContacts()
	assert.NotEqual(t, "Empty! Please, populate the routingtable...", contacts)
}

func TestGetBucketIndex(t *testing.T) {
	adr := address.New("127.0.0.1")
	meId := kademliaid.FromString("ffffffffffffffffffffffffffffffffffffffff")
	me := contact.NewContact(meId, adr)
	rt := routingtable.NewRoutingTable(me)

	// should return the correct bucket index
	target := kademliaid.FromString("fffffffffffffffffffffffffffffffffffffff0")
	bucketIndex := rt.GetBucketIndex(target)
	// target differs on 4 bits so should be in bucket with index 156
	assert.Equal(t, 156, bucketIndex)
	// same id as node should be in the closest bucket
	target = kademliaid.FromString("ffffffffffffffffffffffffffffffffffffffff")
	bucketIndex = rt.GetBucketIndex(target)
	assert.Equal(t, 159, bucketIndex)
}
