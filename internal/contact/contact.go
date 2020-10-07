package contact

import (
	"fmt"
	"kademlia/internal/address"
	"kademlia/internal/kademliaid"
	"sort"
	"strings"
)

// Contact definition
// stores the KademliaID, the ip address and the distance
type Contact struct {
	ID       *kademliaid.KademliaID
	Address  *address.Address
	distance *kademliaid.KademliaID
}

// NewContact returns a new instance of a Contact
func NewContact(id *kademliaid.KademliaID, address *address.Address) Contact {
	return Contact{id, address, nil}
}

// CalcDistance calculates the distance to the target and
// fills the contacts distance field
func (contact *Contact) CalcDistance(target *kademliaid.KademliaID) {
	contact.distance = contact.ID.CalcDistance(target)
}

// Less returns true if contact.distance < otherContact.distance
func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.distance.Less(otherContact.distance)
}

// String returns a simple string representation of a Contact
func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address.String())
}

func (contact *Contact) serialize() string {
	return fmt.Sprintf("%s!%s", contact.ID.String(), contact.Address.String())
}

func (contact *Contact) GetDistance() *kademliaid.KademliaID {
	return contact.distance
}

func SerializeContacts(contacts []Contact) string {
	s := ""
	for i, contact := range contacts {
		if i != len(contacts)-1 {
			s += contact.serialize() + " "
		} else {
			s += contact.serialize()
		}
	}
	return s
}

func Deserialize(s *string) (error, *Contact) {
	fields := strings.Split(*s, "!")
	if len(fields) == 0 || len(fields) == 1 {
		return fmt.Errorf(`Failed to deserialize data string: "%v"`, *s), nil
	}
	adr := address.New(fields[1])
	return nil, &Contact{ID: kademliaid.FromString(fields[0]), Address: adr}
}

// ContactCandidates definition
// stores an array of Contacts
type ContactCandidates struct {
	Contacts []Contact
}

// Append an array of Contacts to the ContactCandidates
func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.Contacts = append(candidates.Contacts, contacts...)
}

// GetContacts returns the first count number of Contacts
func (candidates *ContactCandidates) GetContacts(count int) []Contact {
	return candidates.Contacts[:count]
}

// Sort the Contacts in ContactCandidates
func (candidates *ContactCandidates) Sort() {
	sort.Sort(candidates)
}

// Len returns the length of the ContactCandidates
func (candidates *ContactCandidates) Len() int {
	return len(candidates.Contacts)
}

// Swap the position of the Contacts at i and j
// WARNING does not check if either i or j is within range
func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.Contacts[i], candidates.Contacts[j] = candidates.Contacts[j], candidates.Contacts[i]
}

// Less returns true if the Contact at index i is smaller than
// the Contact at index j
func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.Contacts[i].Less(&candidates.Contacts[j])
}
