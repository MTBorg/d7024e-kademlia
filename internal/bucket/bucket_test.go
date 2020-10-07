package bucket_test

import (
	"container/list"
	"kademlia/internal/address"
	"kademlia/internal/bucket"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBucket(t *testing.T) {
	// should create a new and empty bucket
	b := bucket.NewBucket()
	assert.NotNil(t, b)
	assert.Equal(t, b.Len(), 0)
}

func TestAddContact(t *testing.T) {
	var b *bucket.Bucket
	var bList list.List
	adr := address.New("127.0.0.1")

	// Adding a new contact to a non-full bucket
	// should insert the new contact
	id := kademliaid.NewRandomKademliaID()
	c := contact.NewContact(id, adr)
	b = bucket.NewBucket()
	b.AddContact(c)
	bList = b.GetBucketList()
	assert.Equal(t, bList.Front().Value.(contact.Contact).ID, c.ID)
	assert.Equal(t, bList.Front().Value.(contact.Contact).Address, c.Address)

	// Adding a new contact to a full bucket
	b = bucket.NewBucket()
	for i := 0; i < 20; i++ {
		c = contact.NewContact(kademliaid.NewRandomKademliaID(), adr)
		b.AddContact(c)
	}
	fullBucket := b.GetBucketList()
	assert.Equal(t, fullBucket.Len(), 20)

	// should not add the contact
	newContact := contact.NewContact(kademliaid.NewRandomKademliaID(), adr)
	b.AddContact(newContact)
	newFullBucket := b.GetBucketList()
	assert.Equal(t, newFullBucket.Len(), 20)
	for i, j := newFullBucket.Front(), fullBucket.Front(); i != nil || j != nil; i, j = i.Next(), j.Next() {
		assert.True(t, i.Value.(contact.Contact) == j.Value.(contact.Contact))
	}

	// Adding an already existing contact to the bucket
	// should push the contact to the front of the bucket
	b = bucket.NewBucket()
	var testContact contact.Contact
	for i := 0; i < 20; i++ {
		c = contact.NewContact(kademliaid.NewRandomKademliaID(), adr)
		b.AddContact(c)

		if i == 0 {
			testContact = c
		}
	}
	bList = b.GetBucketList()
	assert.Equal(t, bList.Len(), 20)

	assert.False(t, bList.Front().Value.(contact.Contact) == testContact)
	b.AddContact(testContact)
	bList = b.GetBucketList()
	assert.True(t, bList.Front().Value.(contact.Contact) == testContact)
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
