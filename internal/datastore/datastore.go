package datastore

import (
	"fmt"
	"kademlia/internal/kademliaid"
)

type DataMap = map[kademliaid.KademliaID]string

type DataStore struct {
	data DataMap
}

func New() DataStore {
	return DataStore{make(DataMap)}
}

// Insert a value into the store.
// Uses SHA-1 hash of value as key.
func (d *DataStore) Insert(value string) {
	id := kademliaid.NewKademliaID(&value)
	d.data[id] = value
}

// Gets the value from the store associated with the key.
// Returns an empty string if the key is not found because go is an awful
// language and should never have been invented.
func (d *DataStore) Get(key kademliaid.KademliaID) string {
	return d.data[key]
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
	if len(d.data) != 0 {
		s = "map("
		for key, element := range d.data {
			s = fmt.Sprintf("%s \n %x=%s", s, key, element)
		}
		s += "\n)"
	} else {
		s = "map()"
	}
	return s
}
