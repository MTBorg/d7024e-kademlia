package address_test

import (
	"kademlia/internal/address"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	// should not be nil
	assert.NotNil(t, address.New("123"))

}

func TestString(t *testing.T) {

	// should be equal
	adr := address.New("127.0.0.1:3000")
	assert.Equal(t, adr.String(), "127.0.0.1:")
	adr = address.New("127.0.0.1")
	assert.Equal(t, adr.String(), "127.0.0.1:")

}

func TestGetHost(t *testing.T) {

	// should be the same host address
	adr := address.New("127.0.0.1:1776")
	assert.Equal(t, adr.GetHost(), "127.0.0.1")
}

func TestGetPortAsInt(t *testing.T) {

	adr := address.New("127.0.0.1:3000")
	port, _ := adr.GetPortAsInt()
	assert.Equal(t, port, 0)
}
