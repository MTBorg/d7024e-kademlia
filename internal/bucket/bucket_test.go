package bucket_test

import (
	"kademlia/internal/address"
	"kademlia/internal/bucket"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBucket(t *testing.T) {
	// should create a bucket
	b := bucket.NewBucket()
	assert.NotNil(t, b)
	assert.Equal(t, b.Len(), 0)
}

func TestAddContact(t *testing.T) {
	// shoud be the same contact
	id := kademliaid.NewRandomKademliaID()
	adr := address.New("127.0.0.1")
	c := contact.NewContact(id, adr)
	b := bucket.NewBucket()
	b.AddContact(c)
	contact := b.GetContactAndCalcDistance(id)[0]
	assert.Equal(t, contact.ID, c.ID)
	assert.Equal(t, contact.Address, c.Address)
}

func TestGetContactAndCalcDistance(t *testing.T) {

	// should be the same id
	id1 := kademliaid.FromString("1111111111111111111100000000000000000000")
	id2 := kademliaid.FromString("0000000000000000000011111111111111111111")
	adr := address.New("127.0.0.1")
	c1 := contact.NewContact(id1, adr)
	c2 := contact.NewContact(id2, adr)
	b := bucket.NewBucket()
	b.AddContact(c1)
	b.AddContact(c2)
	c2 = b.GetContactAndCalcDistance(id1)[0]
	c1 = b.GetContactAndCalcDistance(id1)[1]
	assert.Equal(t, c1.ID, kademliaid.FromString("1111111111111111111100000000000000000000"))
	assert.Equal(t, c2.ID, kademliaid.FromString("0000000000000000000011111111111111111111"))

	// should have the correct distance
	assert.Equal(t, c1.GetDistance(), kademliaid.FromString("0000000000000000000000000000000000000000"))
	assert.Equal(t, c2.GetDistance(), kademliaid.FromString("1111111111111111111111111111111111111111"))
}

func TestGetContactAndCalcDistanceNoRequestor(t *testing.T) {

	//should return the contacts
	id1 := kademliaid.FromString("1111111111111111111100000000000000000000")
	id2 := kademliaid.FromString("0000000000000000000011111111111111111111")
	adr := address.New("127.0.0.1")
	c1 := contact.NewContact(id1, adr)
	b := bucket.NewBucket()
	b.AddContact(c1)
	carr := b.GetContactAndCalcDistanceNoRequestor(id1, id2)
	assert.Equal(t, carr[0].ID, id1)
	assert.Equal(t, b.Len(), 1)

	// should not return the requestor
	carr = b.GetContactAndCalcDistanceNoRequestor(id1, id1)
	assert.Nil(t, carr)
}

func TestLen(t *testing.T) {
	// should return 0
	b := bucket.NewBucket()
	assert.Equal(t, b.Len(), 0)

	// shoud return a int 3
	adr := address.New("127.0.0.1")
	c1 := contact.NewContact(kademliaid.NewRandomKademliaID(), adr)
	c2 := contact.NewContact(kademliaid.NewRandomKademliaID(), adr)
	c3 := contact.NewContact(kademliaid.NewRandomKademliaID(), adr)
	b.AddContact(c1)
	b.AddContact(c2)
	b.AddContact(c3)
	assert.Equal(t, b.Len(), 3)
}
