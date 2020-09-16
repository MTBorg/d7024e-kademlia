package pong_test

import (
	"kademlia/internal/rpccommands/pong"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	p := pong.New()
	options := []string{"hello", "abc"}
	//Should never return an error
	assert.NoError(t, p.ParseOptions(&options))
}
