package datastore_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
)

type SenderMock struct {
	mock.Mock
}

func (m *SenderMock) Send(data string, target *address.Address) error {
	args := m.Called(data, target)
	return args.Error(0)
}

func TestForget(t *testing.T) {
	d := datastore.New()
	value := "hello"
	key := kademliaid.NewKademliaID(&value)
	contacts := &[]contact.Contact{}
	d.Insert(value, contacts, nil, nil)

	// should mark the data as forgotten
	err := d.Forget(&key)
	assert.Nil(t, err)

	// should return an error when trying to forget a non-existing value
	value = "byebye"
	key = kademliaid.NewKademliaID(&value)
	err = d.Forget(&key)
	assert.NotNil(t, err)
}

func TestGet(t *testing.T) {
	var d datastore.DataStore

	// Should be able to  get
	d = datastore.New()
	value := "hello"
	contacts := &[]contact.Contact{}
	d.Insert(value, contacts, nil, nil)
	assert.Equal(t, d.Get(kademliaid.NewKademliaID(&value)), "hello")

	// Should not be able to get non-existent key
	d = datastore.New()
	value = "hello"
	assert.Equal(t, d.Get(kademliaid.NewKademliaID(&value)), "")
}

func TestInsert(t *testing.T) {
	var d datastore.DataStore
	var contacts *[]contact.Contact
	value := "hello"
	hash := kademliaid.NewKademliaID(&value)
	originatorId := kademliaid.NewRandomKademliaID()
	originator := contact.NewContact(originatorId, address.New("localhost:3000"))

	//should be able to insert as non originator
	d = datastore.New()
	contacts = &[]contact.Contact{}
	d.Insert(value, contacts, nil, nil)
	assert.Equal(t, d.Get(kademliaid.NewKademliaID(&value)), "hello")

	//should be able to insert as originator
	// d = datastore.New()
	// contacts = &[]contact.Contact{}
	// log.Error().Msgf("HEJSAN %s", os.Getenv("REFRESH_TIME"))
	// d.Insert(value, contacts, &originator, nil)
	// assert.Equal(t, d.Get(kademliaid.NewKademliaID(&value)), "hello")

	//should send refresh RPCs if originator
	d = datastore.New()
	otherContactId := kademliaid.NewRandomKademliaID()
	otherContact := contact.NewContact(otherContactId, address.New("localhost:3000"))
	contacts = &[]contact.Contact{otherContact}
	var senderMock *SenderMock
	senderMock = new(SenderMock)
	senderMock.On("Send", mock.Anything, otherContact.Address).Return(nil)
	d.Insert(value, contacts, &originator, senderMock)
	// Sleep for a bit so that the select case can trigger in the goroutine
	time.Sleep(time.Second * 6) // Sleep for one second more than the default time
	senderMock.AssertExpectations(t)

	// should drop the value when the TTL runs out
	d = datastore.New()
	os.Setenv("REFRESH_TIME", "10")
	os.Setenv("TTL_TIME", "1")
	contacts = &[]contact.Contact{}
	d.Insert(value, contacts, nil, nil)
	// Sleep for a bit so that the select case can trigger in the goroutine
	time.Sleep(time.Second * 2)
	assert.Equal(t, "", d.Get(hash))

	// should stop refreshing and drop the value if the originator has marked
	// the value as forgotten
	d = datastore.New()
	os.Setenv("REFRESH_TIME", "1")
	contacts = &[]contact.Contact{}
	addr := address.New("127.0.0.1:1234")
	me := contact.NewContact(kademliaid.NewRandomKademliaID(), addr)
	d.Insert(value, contacts, &me, nil)
	d.Forget(&hash)
	assert.Equal(t, value, d.Get(hash))
	time.Sleep(time.Millisecond * 1100)
	assert.Equal(t, "", d.Get(hash))
}

func TestEntriesAsString(t *testing.T) {
	var d datastore.DataStore

	//should print map() when empty
	d = datastore.New()
	assert.Equal(t, d.EntriesAsString(), "map()")

	//should print key-value pairs when non-empty
	d = datastore.New()
	v1, v2 := "hello", "world"
	contacts := &[]contact.Contact{}
	d.Insert(v1, contacts, nil, nil)
	d.Insert(v2, contacts, nil, nil)
	whitespaces := regexp.MustCompile(`\s+`)
	fmt.Println(whitespaces.ReplaceAllString(d.EntriesAsString(), ""))
	assert.Contains(t, d.EntriesAsString(), fmt.Sprintf("%x=%s", kademliaid.NewKademliaID(&v1), v1))
	assert.Contains(t, d.EntriesAsString(), fmt.Sprintf("%x=%s", kademliaid.NewKademliaID(&v2), v2))
}

func TestDrop(t *testing.T) {
	var d datastore.DataStore

	d = datastore.New()
	v1, v2 := "hello", "world"
	contacts := &[]contact.Contact{}
	d.Insert(v1, contacts, nil, nil)
	d.Insert(v2, contacts, nil, nil)

	// should delete the entry
	d.Drop("hello")
	assert.Equal(t, "", d.Get(kademliaid.NewKademliaID(&v1)))
	assert.Equal(t, v2, d.Get(kademliaid.NewKademliaID(&v2)))
}
