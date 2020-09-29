package routingtable

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"kademlia/internal/bucket"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
)

const bucketSize = 20

// RoutingTable definition
// keeps a refrence contact of me and an array of buckets
type RoutingTable struct {
	me      contact.Contact
	buckets [kademliaid.IDLength * 8]*bucket.Bucket
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewRoutingTable(me contact.Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < kademliaid.IDLength*8; i++ {
		routingTable.buckets[i] = bucket.NewBucket()
	}
	routingTable.me = me
	return routingTable
}

// GetMe returns the me contact stored in the routing table
func (routingTable *RoutingTable) GetMe() *contact.Contact {
	return &routingTable.me
}

// AddContact add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddContact(contact contact.Contact) {
	if contact.ID.Equals(routingTable.me.ID) != true {
		bucketIndex := routingTable.GetBucketIndex(contact.ID)
		bucket := routingTable.buckets[bucketIndex]
		bucket.AddContact(contact)
	} else {
		log.Warn().Str("me", routingTable.me.String()).Str("contact", contact.String()).Msg("Tried to add self as contact")
	}
}

// FindClosestContacts finds the count closest Contacts to the target in the RoutingTable
func (routingTable *RoutingTable) FindClosestContacts(target *kademliaid.KademliaID, requestorID *kademliaid.KademliaID, count int) []contact.Contact {
	var candidates contact.ContactCandidates
	bucketIndex := routingTable.GetBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	if requestorID != nil {
		candidates.Append(bucket.GetContactAndCalcDistanceNoRequestor(target, requestorID))
	} else {
		candidates.Append(bucket.GetContactAndCalcDistance(target))
	}

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < kademliaid.IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			if requestorID != nil {
				candidates.Append(bucket.GetContactAndCalcDistanceNoRequestor(target, requestorID))
			} else {
				candidates.Append(bucket.GetContactAndCalcDistance(target))
			}
		}
		if bucketIndex+i < kademliaid.IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			if requestorID != nil {
				candidates.Append(bucket.GetContactAndCalcDistanceNoRequestor(target, requestorID))
			} else {
				candidates.Append(bucket.GetContactAndCalcDistance(target))
			}
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.GetContacts(count)
}

// GetContacts returns a string with contacts in the nodes routingtable, the
// contacts are printed with the bucket they reside in. The bucket number
// have been converted so that bucket 0 is the bucket with contacts of inside
// the range 2^0 to 2^1 from the node
func (routingTable *RoutingTable) GetContacts() string {
	if routingTable == nil {
		return "The node is not initilized, it does not contain a routing table or any contacts"
	}
	s := "\nContacts:\n"
	numContacts := 0
	for i, bucket := range routingTable.buckets {
		if bucket != nil && bucket.Len() > 0 {
			convertedBucketIndex := 160 - i
			s += fmt.Sprintf(
				"\nBucket %d:\n",
				convertedBucketIndex,
			)
			contacts := bucket.GetContactAndCalcDistance(routingTable.me.ID)
			numContacts += len(contacts)
			for _, c := range contacts {
				s += fmt.Sprintf("%s\n", c.String())
			}
		}

	}
	s += fmt.Sprintf("\nEnd of contacts.\nTotal number of contacts: %d", numContacts)
	return s
}

// getBucketIndex get the correct Bucket index for the KademliaID
func (routingTable *RoutingTable) GetBucketIndex(id *kademliaid.KademliaID) int {
	distance := id.CalcDistance(routingTable.me.ID)
	for i := 0; i < kademliaid.IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return kademliaid.IDLength*8 - 1
}

func (routingTable *RoutingTable) GetBucket(index int) *bucket.Bucket {
	return routingTable.buckets[index]
}
