package rpcpool_test

import (
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpcpool"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	pool := *rpcpool.New()
	assert.NotNil(t, pool)

}

func TestAdd(t *testing.T) {
	pool := *rpcpool.New()
	id := kademliaid.NewRandomKademliaID()
	assert.Nil(t, pool.GetEntry(id))
	pool.Add(id)
	assert.NotNil(t, pool.GetEntry(id))

}

func TestGetEntry(t *testing.T) {
	pool := *rpcpool.New()
	id := kademliaid.NewRandomKademliaID()
	assert.Nil(t, pool.GetEntry(id))
	pool.Add(id)
	assert.NotNil(t, pool.GetEntry(id))
}

func TestDelete(t *testing.T) {
	pool := *rpcpool.New()
	id := kademliaid.NewRandomKademliaID()
	pool.Add(id)
	assert.NotNil(t, pool.GetEntry(id))
	pool.Delete(id)
	assert.Nil(t, pool.GetEntry(id))
}
