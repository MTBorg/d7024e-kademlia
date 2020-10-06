package datastore

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"kademlia/internal/kademliaid"
	"time"
)

type DataMap = map[kademliaid.KademliaID]Data

type DataStore struct {
	store DataMap
}

type Data struct {
	value   string
	restart chan bool
}

func New() DataStore {
	return DataStore{make(DataMap)}
}

// Insert a data into the store.
// Uses SHA-1 hash of value as key.
// Starts a TTL timer on the data
func (d *DataStore) Insert(value string) {
	id := kademliaid.NewKademliaID(&value)
	data := Data{}
	data.value = value
	data.restart = make(chan bool)
	d.store[id] = data
	d.StartRefreshTimer(data) // If successful insert, we start the TTL

}

// Gets the value from the store associated with the key.
// Returns an empty string if the key is not found because go is an awful
// language and should never have been invented.
func (d *DataStore) Get(key kademliaid.KademliaID) string {
	data := d.store[key]
	d.RestartRefreshTimer(data) // Data is requested
	return data.value
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
			s = fmt.Sprintf("%s \n %x=%s", s, key, element) // TODO: map is now storing the data obj and not a string!
		}
		s += "\n)"
	} else {
		s = "map()"
	}
	return s
}

func (d *DataStore) StartRefreshTimer(data Data) {
	go func() {
		for {
			t := time.Hour
			select {
			case <-data.restart:
				log.Trace().Str("Data", data.value).Msg("Restarted data refresh timer")
			case <-time.After(t):
				log.Trace().Str("Data", data.value).Msg("No refresh done on data, data is silently deleted...")
				go d.Drop(data.value)
			}
		}
	}()
}

func (d *DataStore) RestartRefreshTimer(data Data) {
	data.restart <- true // restart the refresh timer
}
