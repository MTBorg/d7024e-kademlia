package kademliaid_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"kademlia/internal/kademliaid"
)

func TestCalcDistance(t *testing.T) {
	kadID := kademliaid.KademliaID{}

	// should be equal because ids are the same
	id := kademliaid.NewRandomKademliaID()
	id2 := id
	assert.Equal(t, &kadID, id.CalcDistance(id2))

}

func TestLess(t *testing.T) {

	// should return false
	id := kademliaid.NewRandomKademliaID()
	id2 := id
	assert.False(t, id.Less(id2))

}

func TestNewKademliaID(t *testing.T) {
	// should return a hash
	data := "TEST"
	hash1 := kademliaid.NewKademliaID(&data)
	assert.NotNil(t, hash1)

	// should return the same hash
	assert.Equal(t, kademliaid.NewKademliaID(&data), kademliaid.NewKademliaID(&data))
}

func TestFromString(t *testing.T) {
	// should be the same
	s := kademliaid.NewRandomKademliaID().String()
	id := kademliaid.FromString(s).String()
	assert.Equal(t, s, id)
}

func TestString(t *testing.T) {
	// should be the same
	id := kademliaid.NewRandomKademliaID()
	s := id.String()
	assert.Equal(t, id, kademliaid.FromString(s))

}

func TestEquals(t *testing.T) {
	// should return true since the both ids are the same
	id := kademliaid.FromString("0000000000FFFFFFFFFF")
	id2 := kademliaid.FromString("0000000000FFFFFFFFFF")
	assert.True(t, id.Equals(id2))

	// should return false since we have generated a new id for id2
	id2 = kademliaid.NewRandomKademliaID()
	assert.False(t, id.Equals(id2))

}

func TestNewRandomKademliaID(t *testing.T) {
	// should return different ids
	id := kademliaid.NewRandomKademliaID()
	assert.False(t, id.Equals(kademliaid.NewRandomKademliaID()))

}
