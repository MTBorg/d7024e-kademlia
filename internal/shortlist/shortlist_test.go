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

	// should not be marked as returned value
	for i := 0; i < k; i++ {
		assert.False(t, sl.Entries[i].ReturnedValue)
	}
}

func TestLess(t *testing.T) {
	var sl shortlist.Shortlist

	// Should return true if the second element is nil
	sl = shortlist.Shortlist{Entries: [5]*shortlist.Entry{&shortlist.Entry{}, nil}}
	assert.Equal(t, true, sl.Less(0, 1))

	// Should return false if the first element is nil and second is not
	sl = shortlist.Shortlist{Entries: [5]*shortlist.Entry{nil, &shortlist.Entry{}}}
	assert.Equal(t, false, sl.Less(0, 1))
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

	// Should be able to sort if list is not full
	contacts = []contact.Contact{}
	target = kademliaid.FromString(strings.Repeat("0", 40))
	for i := 0; i < 1; i++ {
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

func TestGetContacts(t *testing.T) {
	adr := address.New("127.0.0.1")
	contacts := []contact.Contact{}
	target := kademliaid.FromString(strings.Repeat("0", 40))
	for i := 0; i < 4; i++ {
		s := strconv.Itoa(i+1) + strings.Repeat("0", 39) // +1 to avoid all zeroes (equal to target)
		contact := contact.NewContact(kademliaid.FromString(s), adr)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}
	sl := shortlist.NewShortlist(target, contacts)

	cs := sl.GetContacts()
	assert.Equal(t, 4, len(cs))
}

func TestDrop(t *testing.T) {
	// should remove a contact
	adr := address.New("127.0.0.1")
	contacts := []contact.Contact{}
	target := kademliaid.FromString(strings.Repeat("0", 40))
	for i := 0; i < 4; i++ {
		s := strconv.Itoa(i+1) + strings.Repeat("0", 39) // +1 to avoid all zeroes (equal to target)
		contact := contact.NewContact(kademliaid.FromString(s), adr)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}
	sl := shortlist.NewShortlist(target, contacts)
	assert.Equal(t, 4, sl.Len())
	contact := sl.Entries[2].Contact
	assert.NotNil(t, contact)
	sl.Drop(&contact)
	assert.Equal(t, 3, sl.Len())
	for i := 0; i < len(sl.Entries); i++ {
		if sl.Entries[i] != nil {
			assert.NotEqual(t, contact, sl.Entries[i].Contact)
		}
	}
}
