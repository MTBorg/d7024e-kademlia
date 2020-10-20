package datastore

import (
	"errors"
	"fmt"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpc"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type DataMap = map[kademliaid.KademliaID]*Data

type DataStore struct {
	store DataMap
}

type Data struct {
	value      string
	restart    chan bool
	Contacts   *[]contact.Contact
	originator bool
	forgotten  bool
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
func (d *DataStore) Insert(value string, contacts *[]contact.Contact, originator *contact.Contact, sender rpc.Sender) {
	id := kademliaid.NewKademliaID(&value)
	data := Data{}
	data.value = value
	data.restart = make(chan bool)
	data.Contacts = contacts
	data.originator = originator != nil
	d.store[id] = &data
	data.forgotten = false

	// If this is the node that the data was originally stored at
	// then we want to refresh rather than delete it
	if originator != nil {
		d.StartRefreshTimer(data, originator, sender) // If successful insert, we start the TTL
	} else {
		d.StartTTLTimer(data, sender)
	}
}

// Gets the value from the store associated with the key.
// Returns an empty string if the key is not found because go is an awful
// language and should never have been invented.
func (d *DataStore) Get(key kademliaid.KademliaID) string {
	data := d.store[key]
	if data != nil {
		if !data.originator {
			d.RestartRefreshTimer(*data)
		}
		return data.value

	}
	return ""

}

// Drop removes the data from the store
func (d *DataStore) Drop(value string) {
	id := kademliaid.NewKademliaID(&value)
	delete(d.store, id)
}

func (d *DataStore) Forget(key *kademliaid.KademliaID) error {
	if d.store[*key] != nil {
		log.Trace().Str("Hash", key.String()).Msg("Marking entry as forgotten")
		d.store[*key].forgotten = true
		return nil
	} else {
		log.Trace().Str("Key", key.String()).Msg("Tried to forget non-existent entry")
		return errors.New("Tried to forget non-existent entry")
	}
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

func (d *DataStore) StartTTLTimer(data Data, sender rpc.Sender) {
	go func() {
		ttlTime, err := strconv.Atoi(os.Getenv("TTL_TIME"))
		if err != nil {
			log.Error().Msgf("Failed to convert env variable TTL_TIME from string to int: %s", err)
			ttlTime = 10
		}
		for {
			ttlTimer := time.Duration(float64(time.Second) * float64(ttlTime))
			select {
			case <-data.restart:
				log.Trace().Str("Data", data.value).Msg("Restarted data refresh timer")
			case <-time.After(ttlTimer):
				log.Trace().Str("Data", data.value).Msg("No refresh done on data, data is silently deleted...")
				d.Drop(data.value)
				return
			}
		}
	}()
}

func (d *DataStore) StartRefreshTimer(data Data, originator *contact.Contact, sender rpc.Sender) {
	go func() {
		for {
			refreshTime, err := strconv.Atoi(os.Getenv("REFRESH_TIME"))
			if err != nil {
				log.Error().Msgf("Failed to convert env variable REFRESH_TIME from string to int: %s", err)
				refreshTime = 5
			}
			refreshTimer := time.Duration(float64(time.Second) * float64(refreshTime))

			select {
			case <-time.After(refreshTimer):
				hash := kademliaid.NewKademliaID(&data.value)

				//If the entry has been marked as forgotten, stop refreshing
				if d.store[hash].forgotten {
					log.Trace().
						Str("Hash", hash.String()).
						Msg("Entry has been marked as forgotten, deleting entry and stopping refresh")
					delete(d.store, hash)
					return
				}

				log.Trace().Str("Hash", hash.String()).Msg("Sending refreshes")
				for _, contact := range *data.Contacts {
					refresh := rpc.New(originator.ID,
						"REFRESH "+hash.String(), contact.Address)
					refresh.Send(sender, refresh.Target)
				}
			}
		}
	}()
}

func (d *DataStore) RestartRefreshTimer(data Data) {
	data.restart <- true // restart the refresh timer
}
