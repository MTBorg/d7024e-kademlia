package datastore_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
)

func TestGet(t *testing.T) {
	var d datastore.DataStore

	// Should be able to  get
	d = datastore.New()
	value := "hello"
	d.Insert(value)
	assert.Equal(t, d.Get(kademliaid.NewKademliaID(&value)), "hello")

	// Should not be able to get non-existent key
	d = datastore.New()
	value = "hello"
	assert.Equal(t, d.Get(kademliaid.NewKademliaID(&value)), "")
}

func TestInsert(t *testing.T) {
	var d datastore.DataStore

	//should be able to insert
	d = datastore.New()
	value := "hello"
	d.Insert(value)
	assert.Equal(t, d.Get(kademliaid.NewKademliaID(&value)), "hello")
}

func TestEntriesAsString(t *testing.T) {
	var d datastore.DataStore

	//should print map() when empty
	d = datastore.New()
	assert.Equal(t, d.EntriesAsString(), "map()")

	//should print key-value pairs when non-empty
	d = datastore.New()
	v1, v2 := "hello", "world"
	d.Insert(v1)
	d.Insert(v2)
	whitespaces := regexp.MustCompile(`\s+`)
	fmt.Println(whitespaces.ReplaceAllString(d.EntriesAsString(), ""))
	assert.Contains(t, d.EntriesAsString(), fmt.Sprintf("%x=%s", kademliaid.NewKademliaID(&v1), v1))
	assert.Contains(t, d.EntriesAsString(), fmt.Sprintf("%x=%s", kademliaid.NewKademliaID(&v2), v2))
}

func TestDrop(t *testing.T) {
	var d datastore.DataStore

	d = datastore.New()
	v1, v2 := "hello", "world"
	d.Insert(v1)
	d.Insert(v2)

	// should delete the entry
	d.Drop("hello")
	assert.Equal(t, "", d.Get(kademliaid.NewKademliaID(&v1)))
	assert.Equal(t, v2, d.Get(kademliaid.NewKademliaID(&v2)))
}
