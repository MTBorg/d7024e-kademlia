package address_test

import (
	"kademlia/internal/address"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Valid Address with IP and Port
	// should return a new address with the same IP is input
	inAddr := "127.1.1.1:1234"
	outAddr := address.New(inAddr)
	assert.Equal(t, "127.1.1.1", outAddr.GetHost())
	outPort, _ := outAddr.GetPortAsInt()
	assert.Equal(t, 1776, outPort)

	// Address with only IP
	// should still return a new address with the default port
	inAddr = "127.1.1.1"
	outAddr = address.New(inAddr)
	assert.Equal(t, "127.1.1.1", outAddr.GetHost())
	outPort, _ = outAddr.GetPortAsInt()
	assert.Equal(t, 1776, outPort)

	// Invalid address
	// should return an empty address
	outAddr = address.New("hajhaj")
	assert.Equal(t, &address.Address{}, outAddr)
}

func TestString(t *testing.T) {
	// should return the address as a string
	inAddr := "127.0.0.1:1776"
	adr := address.New(inAddr)
	assert.Equal(t, adr.String(), inAddr)
}

func TestGetHost(t *testing.T) {
	// should be the same host address
	adr := address.New("127.0.0.1:1776")
	assert.Equal(t, adr.GetHost(), "127.0.0.1")
}

func TestGetPortAsInt(t *testing.T) {
	// should return the port as int
	adr := address.New("127.0.0.1:1776")
	port, err := adr.GetPortAsInt()
	assert.Nil(t, err)
	assert.Equal(t, port, 1776)
}
