package bucket

import (
	"container/list"
	. "kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	. "kademlia/internal/kademliaid"
	"kademlia/internal/network"
	"kademlia/internal/rpc"
	"kademlia/internal/rpcpool"
	"time"
)

// Bucket definition
// contains a List of contacts, the ID of the node and an RPCPool used when
// a contact has to be pinged to determine if it's alive (responed with pong).
type Bucket struct {
	list    *list.List
	nodeID  *kademliaid.KademliaID
	rpcPool *rpcpool.RPCPool
}

const tWaitForPong = 5
const bucketSize = 20

// NewBucket returns a new instance of a Bucket
func NewBucket(nodeID *kademliaid.KademliaID) *Bucket {
	bucket := &Bucket{}
	bucket.list = list.New()
	bucket.nodeID = nodeID
	bucket.rpcPool = rpcpool.New()
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *Bucket) AddContact(contact Contact) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
		} else {
			// send PING to least recently seen (LRS) contact
			LRSContact := bucket.list.Back()
			rpc := rpc.New(bucket.nodeID, "PING", LRSContact.Value.(Contact).Address)
			network.Net.SendPingMessage(&rpc)
			bucket.rpcPool.Add(rpc.RPCId)

			// wait for at most tWaitForPong and if the LRS contact doesn't respond
			// evict it from the bucket and insert the new, active, contact
			res := bucket.rpcPool.GetEntry(rpc.RPCId)
			select {
			case <-res.Channel:
				bucket.list.PushFront(LRSContact)
			case <-time.After(tWaitForPong * time.Second):
				bucket.list.Remove(LRSContact)
				bucket.list.PushFront(contact)
			}
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *Bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// GetContactAndCalcDistance returns an array of Contacts where the distance
// has already been calculated. This array will never contain a contact with
// the same nodeID as the requestorID.
func (bucket *Bucket) GetContactAndCalcDistanceNoRequestor(target *KademliaID, requestorID *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		if !contact.ID.Equals(requestorID) {
			contact.CalcDistance(target)
			contacts = append(contacts, contact)
		}
	}
	return contacts
}

// Len return the size of the bucket
func (bucket *Bucket) Len() int {
	return bucket.list.Len()
}
