package put_test

import (
	"github.com/stretchr/testify/assert"
	"kademlia/internal/commands/put"
	"testing"
)

func TestParseOptions(t *testing.T) {
	var putCmd *put.Put
	var err error

	// should not return an error if content is specified
	putCmd = new(put.Put)
	err = putCmd.ParseOptions([]string{"address", "content"})
	assert.Nil(t, err)

	// should return an error
	putCmd = new(put.Put)
	err = putCmd.ParseOptions([]string{})
	assert.NotNil(t, err)

}
