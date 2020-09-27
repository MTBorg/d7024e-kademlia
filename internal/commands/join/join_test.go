package join_test

import (
	"kademlia/internal/commands/join"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	// TODO: Not tested since .net
}

func TestParseOptions(t *testing.T) {
	var joinCmd *join.Join
	assert.Nil(t, joinCmd.ParseOptions([]string{}))
	assert.Nil(t, joinCmd.ParseOptions([]string{"test", "test"}))
}

func TestPrintUsage(t *testing.T) {
	// should be equal
	var joinCmd *join.Join
	assert.Equal(t, joinCmd.PrintUsage(), "Usage: join")

}
