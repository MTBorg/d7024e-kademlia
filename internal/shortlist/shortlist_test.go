package shortlist_test

import (
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/shortlist"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: USE ENV CONST
const k = 5

func TestNewShortlist(t *testing.T) {

	addr := address.New("address")
	target := kademliaid.NewRandomKademliaID()
	candidates := [k]contact.Contact{
		contact.NewContact(kademliaid.NewRandomKademliaID(), addr),
		contact.NewContact(kademliaid.NewRandomKademliaID(), addr),
		contact.NewContact(kademliaid.NewRandomKademliaID(), addr),
		contact.NewContact(kademliaid.NewRandomKademliaID(), addr),
		contact.NewContact(kademliaid.NewRandomKademliaID(), addr),
	}
	sl := shortlist.NewShortlist(target, candidates[:])

	// should have k entries
	assert.Len(t, sl.Entries, k)

	// should not be marked as probed
	for i := 0; i < k; i++ {
		assert.False(t, sl.Entries[i].Probed)
	}

	// should not be marked as active
	for i := 0; i < k; i++ {
		assert.False(t, sl.Entries[i].Active)
	}
}

func TestLen(t *testing.T) {
	addr := address.New("address")
	candidates := []contact.Contact{
		contact.NewContact(kademliaid.NewRandomKademliaID(), addr),
		contact.NewContact(kademliaid.NewRandomKademliaID(), addr),
		contact.NewContact(kademliaid.NewRandomKademliaID(), addr),
	}
	target := kademliaid.NewRandomKademliaID()
	sl := shortlist.NewShortlist(target, candidates[:])

	assert.Equal(t, 3, sl.Len())
}

func TestAdd(t *testing.T) {
	adr := address.New("localhost")
	targetID := kademliaid.NewRandomKademliaID()
	var sl *shortlist.Shortlist

	// Should not add a contact if already exists
	c := contact.NewContact(kademliaid.NewRandomKademliaID(), adr)
	c.CalcDistance(targetID)
	sl = shortlist.NewShortlist(targetID, []contact.Contact{c})
	sl.Add(&c)
	assert.Equal(t, 1, sl.Len())

	// Should add a contact if it doesn't exist
	c2 := contact.NewContact(kademliaid.NewRandomKademliaID(), adr)
	c2.CalcDistance(targetID)
	sl = shortlist.NewShortlist(targetID, []contact.Contact{c2})
	sl.Add(&c)
	assert.Equal(t, 2, sl.Len())

	// Should add and sort if list is not full
	contacts := []contact.Contact{}
	target := kademliaid.FromString(strings.Repeat("0", 40))
	for i := 0; i < 4; i++ {
		s := strconv.Itoa(i+1) + strings.Repeat("0", 39) // +1 to avoid all zeroes (equal to target)
		contact := contact.NewContact(kademliaid.FromString(s), adr)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}
	sl = shortlist.NewShortlist(target, contacts)
	c = contact.NewContact(kademliaid.FromString(strings.Repeat("0", 39)+"1"), adr)
	c.CalcDistance(target)
	sl.Add(&c)
	assert.Equal(t, c, sl.Entries[0].Contact)

	// Should add and sort if list is full
	contacts = []contact.Contact{}
	target = kademliaid.FromString(strings.Repeat("0", 40))
	for i := 0; i < 5; i++ {
		s := strconv.Itoa(i+1) + strings.Repeat("0", 39) // +1 to avoid all zeroes (equal to target)
		contact := contact.NewContact(kademliaid.FromString(s), adr)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}
	sl = shortlist.NewShortlist(target, contacts)

	c = contact.NewContact(kademliaid.FromString(strings.Repeat("0", 39)+"1"), adr)
	c.CalcDistance(target)
	sl.Add(&c)
	assert.Equal(t, c, sl.Entries[0].Contact)
}
