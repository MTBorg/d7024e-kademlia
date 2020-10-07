package contact_test

import (
	"fmt"
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
	id1 := kademliaid.NewRandomKademliaID()
	id2 := kademliaid.NewRandomKademliaID()
	adr := address.New("127.0.0.1")
	c := contact.NewContact(id1, adr)

	// should update the distance of the contact
	assert.Nil(t, c.GetDistance())
	c.CalcDistance(id2)
	assert.NotNil(t, c.GetDistance())
}

func TestLess(t *testing.T) {
	adr := address.New("127.0.0.1")

	// should return true if the contact is closer than the other contact
	targetID := kademliaid.FromString("ffffffffffffffffffffffffffffffffffffffff")
	id1 := kademliaid.FromString("00000000000000000000000000000000000000ff")
	id2 := kademliaid.FromString("000000000000000000000000000000000000000f")
	c1 := contact.NewContact(id1, adr)
	c2 := contact.NewContact(id2, adr)
	c1.CalcDistance(targetID)
	c2.CalcDistance(targetID)
	assert.True(t, c1.Less(&c2))

	// should return false if the contact is further away than the other contact
	assert.False(t, c2.Less(&c1))
}

func TestString(t *testing.T) {
	// should return a string with contact
	id := kademliaid.FromString("1111111111111111111100000000000000000000")
	adr := address.New("127.0.0.1")
	c := contact.NewContact(id, adr)
	assert.Equal(t, c.String(), "contact(\"1111111111111111111100000000000000000000\", \"127.0.0.1:1776\")")
}

func TestSerializeContacts(t *testing.T) {
	addr := address.New("127.0.0.1")
	id1 := kademliaid.NewRandomKademliaID()
	id2 := kademliaid.NewRandomKademliaID()
	cs := []contact.Contact{contact.NewContact(id1, addr), contact.NewContact(id2, addr)}

	// should return a string representation of the contacts
	csStr := contact.SerializeContacts(cs)
	expected := fmt.Sprintf("%s!%s %s!%s", cs[0].ID.String(), cs[0].Address.String(), cs[1].ID.String(), cs[1].Address.String())
	assert.Equal(t, expected, csStr)
}

func TestDeserialize(t *testing.T) {
	addr := address.New("127.0.0.1")
	id1 := kademliaid.NewRandomKademliaID()

	// Valid contact in string format
	// should return the contact
	input := fmt.Sprintf("%s!%s", id1.String(), addr.String())
	err, res := contact.Deserialize(&input)
	assert.Nil(t, err)
	assert.Equal(t, id1, res.ID)
	assert.Equal(t, addr, res.Address)

	// Invalid string input
	// should return an error
	input = "asjdhasjd"
	err, res = contact.Deserialize(&input)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

// ---- CONTACT CANDIDATES ----

func TestAppend(t *testing.T) {
	addr := address.New("127.0.0.1")
	candidates := contact.ContactCandidates{}

	// should append all the contacts
	assert.Equal(t, 0, candidates.Len())
	cs := []contact.Contact{contact.NewContact(kademliaid.NewRandomKademliaID(), addr), contact.NewContact(kademliaid.NewRandomKademliaID(), addr)}
	candidates.Append(cs)
	assert.Equal(t, len(cs), len(candidates.Contacts))
}

func TestGetContacts(t *testing.T) {
	addr := address.New("127.0.0.1")
	candidates := contact.ContactCandidates{}
	cs := []contact.Contact{contact.NewContact(kademliaid.NewRandomKademliaID(), addr), contact.NewContact(kademliaid.NewRandomKademliaID(), addr), contact.NewContact(kademliaid.NewRandomKademliaID(), addr)}
	candidates.Append(cs)

	// should return the specified number of contacts
	assert.Equal(t, 3, len(candidates.Contacts))
	gc := candidates.GetContacts(2)
	assert.Equal(t, 2, len(gc))
}

func TestSort(t *testing.T) {
	addr := address.New("127.0.0.1")
	candidates := contact.ContactCandidates{}
	targetID := kademliaid.FromString("ffffffffffffffffffffffffffffffffffffffff")
	id1 := kademliaid.FromString("fffffffffffffffffffffffffffffffffffff000")
	id2 := kademliaid.FromString("ffffffffffffffffffffffffffffffffffffff00")
	id3 := kademliaid.FromString("fffffffffffffffffffffffffffffffffffffff0")
	cs := []contact.Contact{contact.NewContact(id1, addr), contact.NewContact(id2, addr), contact.NewContact(id3, addr)}
	for i := 0; i < len(cs); i++ {
		cs[i].CalcDistance(targetID)
	}
	candidates.Append(cs)

	// should sort the candidates by distance to the target
	assert.False(t, candidates.Contacts[0].Less(&candidates.Contacts[1]))
	assert.False(t, candidates.Contacts[1].Less(&candidates.Contacts[2]))
	candidates.Sort()
	assert.True(t, candidates.Contacts[0].Less(&candidates.Contacts[1]))
	assert.True(t, candidates.Contacts[1].Less(&candidates.Contacts[2]))
}
