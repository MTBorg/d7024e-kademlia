package shortlist

import (
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"sort"
)

// TODO: Figure out a way to set this using env var. This is not as easy since
// the value is used as the array size.
const k = 5 // k-closest

// Entry represents an entry in the shortlist used in the lookup algorithm.
// an entry in the shortlist can be probed (if the RPC has been sent) and
// active (if the probed contact has responded).
type Entry struct {
	Contact       contact.Contact
	Probed        bool
	Active        bool
	ReturnedValue bool
}

// The shortlist used in the lookup algorithm. The entries in the shortlist
// are the k nodes which are closest to the key being searched for.
type Shortlist struct {
	Entries [k]*Entry
	Closest *contact.Contact
	target  *kademliaid.KademliaID
}

func (sl *Shortlist) Swap(i, j int) {
	sl.Entries[i], sl.Entries[j] = sl.Entries[j], sl.Entries[i]
}

func (sl *Shortlist) Less(i, j int) bool {
	if sl.Entries[j] == nil {
		return true
	}
	if sl.Entries[i] == nil {
		return false
	}
	return sl.Entries[i].Contact.Less(&sl.Entries[j].Contact)
}

// Len returns the number of non-null entries in the shortlist
func (shortlist *Shortlist) Len() int {
	length := 0
	for _, entry := range shortlist.Entries {
		if entry != nil {
			length++
		}
	}
	return length
}

// NewShortlist returns a shortlist of k entires given an initial list of
// k (or less) candidates.
//
// NewShortlist assumes that the supplied candidates are sorted such that
// the first candidate is the closest to the key
func NewShortlist(target *kademliaid.KademliaID, candidates []contact.Contact) *Shortlist {
	shortlist := &Shortlist{}
	shortlist.Closest = &candidates[0]
	shortlist.target = target
	for i, contact := range candidates {
		shortlist.Entries[i] = &Entry{contact, false, false, false}
	}
	return shortlist
}

func (sl *Shortlist) GetContacts() []contact.Contact {
	contacts := []contact.Contact{}
	for _, entry := range sl.Entries {
		if entry != nil {
			contacts = append(contacts, entry.Contact)
		}
	}
	return contacts
}

// Add a contact to the shortlist
//
// The contact is only added if it does not already exist in the list and, if
// the list already contains K contacts, only if it's distance to the id is
// less than that of the furthest away node in the current shortlist.
//
// Assumes the shortlist to be sorted.
func (sl *Shortlist) Add(c *contact.Contact) {
	// Check if contact already exists
	for _, entry := range sl.Entries {
		if entry != nil {
			if entry.Contact.ID.Equals(c.ID) {
				return
			}
		}
	}

	// calc distance of new candidate
	c.CalcDistance(sl.target)

	if sl.Len() == k {
		if c.Less(&sl.Entries[k-1].Contact) {
			sl.Entries[k-1] = &Entry{Contact: *c, Active: false, Probed: false}
		}
	} else {
		for i := 0; i < len(sl.Entries); i++ {
			if sl.Entries[i] == nil {
				sl.Entries[i] = &Entry{Contact: *c, Active: false, Probed: false}
				break
			}
		}
	}

	sort.Sort(sl)
	sl.Closest = &sl.Entries[0].Contact
}

// Drop removes a contact from shortlist by setting it to nil
func (sl *Shortlist) Drop(c *contact.Contact) {
	for i := 0; i < len(sl.Entries); i++ {
		if sl.Entries[i] != nil {
			if sl.Entries[i].Contact.ID.Equals(c.ID) {
				sl.Entries[i] = nil
			}
		}
	}
}
