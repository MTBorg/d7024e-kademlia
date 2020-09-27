package routingtable

import (
	"fmt"
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

// AddContact add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddContact(contact contact.Contact) {
	if contact.ID.Equals(routingTable.me.ID) != true {
		bucketIndex := routingTable.GetBucketIndex(contact.ID)
		bucket := routingTable.buckets[bucketIndex]
		bucket.AddContact(contact)
	}

}

// FindClosestContacts finds the count closest Contacts to the target in the RoutingTable
func (routingTable *RoutingTable) FindClosestContacts(target *kademliaid.KademliaID, count int) []contact.Contact {
	var candidates contact.ContactCandidates
	bucketIndex := routingTable.GetBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < kademliaid.IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < kademliaid.IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.GetContacts(count)
}

// GetContacts returns a newline string with contacts in the nodes routingtable
func (routingTable *RoutingTable) GetContacts() string {
	s := ""
	if routingTable == nil {
		return "Empty! Please, populate the routingtable..."
	}
	for _, bucket := range routingTable.buckets {
		if bucket != nil && bucket.Len() > 0 {
			s += fmt.Sprintf("%+v\n", bucket.GetContactAndCalcDistance(routingTable.me.ID))

		}
	}

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
