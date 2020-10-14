package datastore

import (
	"fmt"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpc"
	"kademlia/internal/udpsender"
	"time"

	"github.com/rs/zerolog/log"
)

type DataMap = map[kademliaid.KademliaID]*Data

type DataStore struct {
	store DataMap
}

type Data struct {
	value    string
	restart  chan bool
	Contacts *[]contact.Contact
}

func New() DataStore {
	return DataStore{make(DataMap)}
}

// Insert a data into the store.
// Uses SHA-1 hash of value as key.
// Starts a TTL timer on the data
//
// The originator parameter should point to the node's routingtable's me value,
// i.e. this node's contact representation. If originator is nil, then that
// means that the storage of the value was not initiated on this node, and thus
// this node should not send REFRESH RPCs to the other nodes.
func (d *DataStore) Insert(value string, contacts *[]contact.Contact, originator *contact.Contact, sender *udpsender.UDPSender) {
	id := kademliaid.NewKademliaID(&value)
	data := Data{}
	data.value = value
	data.restart = make(chan bool)
	data.Contacts = contacts
	d.store[id] = &data
	d.StartRefreshTimer(data, originator, sender) // If successful insert, we start the TTL
}

// Gets the value from the store associated with the key.
// Returns an empty string if the key is not found because go is an awful
// language and should never have been invented.
func (d *DataStore) Get(key kademliaid.KademliaID) string {
	data := d.store[key]
	if data != nil {
		d.RestartRefreshTimer(*data)
		return data.value

	}
	return ""

}

// Drop removes the data from the store
func (d *DataStore) Drop(value string) {
	id := kademliaid.NewKademliaID(&value)
	delete(d.store, id)
}

// Pretty printing of store
func (d *DataStore) EntriesAsString() string {
	/* Format as
	map(
		key1=val1
		key2=val2
		...
	)
	if non-empty, otherwise:
	map()
	*/
	var s string
	if len(d.store) != 0 {
		s = "map("
		for key, element := range d.store {
			s = fmt.Sprintf("%s \n %x=%s", s, key, element.value)
		}
		s += "\n)"
	} else {
		s = "map()"
	}
	return s
}

func (d *DataStore) StartRefreshTimer(data Data, originator *contact.Contact, sender *udpsender.UDPSender) {
	go func() {
		for {
			refreshTime := time.Second * 5
			// t := time.Hour
			if originator != nil {
				// If this is the node that the data was originally stored at
				// then we want to refresh rather than delete it
				select {
				case <-time.After(refreshTime):
					hash := kademliaid.NewKademliaID(&data.value)
					log.Trace().Str("Hash", hash.String()).Msg("Sending refreshes")
					for _, contact := range *data.Contacts {
						refresh := rpc.New(originator.ID,
							"REFRESH "+hash.String(), contact.Address)
						refresh.Send(sender, refresh.Target)
					}
				}
			} else {
				t := time.Second * 10
				select {
				case <-data.restart:
					log.Trace().Str("Data", data.value).Msg("Restarted data refresh timer")
				case <-time.After(t):
					log.Trace().Str("Data", data.value).Msg("No refresh done on data, data is silently deleted...")
					d.Drop(data.value)
					return
				}
			}
		}
	}()
}

func (d *DataStore) RestartRefreshTimer(data Data) {
	data.restart <- true // restart the refresh timer
}
