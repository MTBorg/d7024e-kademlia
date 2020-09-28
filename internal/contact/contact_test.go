package contact_test

import (
	"kademlia/internal/address"
	"kademlia/internal/kademliaid"
	"testing"

	"github.com/stretchr/testify/assert"

	"kademlia/internal/contact"
)

func TestNewContact(t *testing.T) {
	// should return a contact
	id := kademliaid.NewRandomKademliaID()
	adr := address.New("127.0.0.1")
	c := contact.NewContact(id, adr)
	assert.NotNil(t, c)
	assert.Equal(t, c.ID, id)
	assert.Equal(t, c.Address, adr)
}

func TestCalcDistance(t *testing.T) {

	// should be equal
	id1 := kademliaid.FromString("1111111111111111111100000000000000000000")
	id2 := kademliaid.FromString("0000000000000000000011111111111111111111")
	adr := address.New("127.0.0.1")
	c := contact.NewContact(id1, adr)
	c.CalcDistance(id2)
	assert.Equal(t, c.GetDistance(), kademliaid.FromString("1111111111111111111111111111111111111111"))
}

func TestLess(t *testing.T) {
	// should be false
	id1 := kademliaid.FromString("1111111111111111111100000000000000000000")
	id2 := kademliaid.FromString("0000000000000000000011111111111111111111")
	adr := address.New("127.0.0.1")
	c1 := contact.NewContact(id1, adr)
	c2 := contact.NewContact(id1, adr)
	c1.CalcDistance(id2)
	c2.CalcDistance(id2)
	assert.False(t, c1.Less(&c2))
	assert.False(t, c2.Less(&c1))

}

func TestString(t *testing.T) {
	// should return a string with contact
	id := kademliaid.FromString("1111111111111111111100000000000000000000")
	adr := address.New("127.0.0.1")
	c := contact.NewContact(id, adr)
	assert.Equal(t, c.String(), "contact(\"1111111111111111111100000000000000000000\", \"127.0.0.1:1776\")")
}
