package rpcpool_test

import (
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpcpool"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock of the HOF supplied to the WithLock function
type FnMock struct {
	mock.Mock
}

func (m *FnMock) Exec() {
	m.Called()
	return
}

func TestNew(t *testing.T) {
	pool := *rpcpool.New()
	assert.NotNil(t, pool)
}

func TestAdd(t *testing.T) {
	pool := *rpcpool.New()
	id := kademliaid.NewRandomKademliaID()
	contact := contact.Contact{ID: id}

	// should add a new entry
	assert.Nil(t, pool.GetEntry(id))
	pool.Add(id, &contact)
	assert.NotNil(t, pool.GetEntry(id))
	assert.NotNil(t, pool.GetEntry(id).Contact)
}

func TestGetEntry(t *testing.T) {
	pool := *rpcpool.New()
	id := kademliaid.NewRandomKademliaID()
	contact := contact.Contact{ID: id}

	// should return the entry
	assert.Nil(t, pool.GetEntry(id))
	pool.Add(id, &contact)
	assert.NotNil(t, pool.GetEntry(id))

}

func TestDelete(t *testing.T) {
	pool := *rpcpool.New()
	id := kademliaid.NewRandomKademliaID()
	contact := contact.Contact{ID: id}

	// should delete the entry
	pool.Add(id, &contact)
	assert.NotNil(t, pool.GetEntry(id))
	pool.Delete(id)
	assert.Nil(t, pool.GetEntry(id))
}

func TestWithLock(t *testing.T) {
	pool := *rpcpool.New()
	funcMock := new(FnMock)
	funcMock.On("Exec").Return()

	// should call the supplied function after locking the pool
	pool.WithLock(func() { funcMock.Exec() })
	funcMock.AssertNumberOfCalls(t, "Exec", 1)
}
