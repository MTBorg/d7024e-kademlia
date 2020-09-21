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
	adr := address.New("127.0.0.1:1776")
	assert.Equal(t, adr.String(), "127.0.0.1:1776")

}
