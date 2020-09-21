package findvalue_test

import (
	"kademlia/internal/rpccommands/findvalue"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	var err error
	find := findvalue.FindValue{}

	//Should return an error if id was not specified
	err = find.ParseOptions(&[]string{})
	assert.EqualError(t, err, "Missing hash")

	//Should not return an error if id was not specified
	err = find.ParseOptions(&[]string{"someid"})
	assert.NoError(t, err)
}
