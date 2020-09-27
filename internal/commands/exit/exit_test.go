package exit_test

import (
	"kademlia/internal/commands/exit"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOptions(t *testing.T) {
	// should be nil
	var exitCmd *exit.Exit
	assert.Nil(t, exitCmd.ParseOptions([]string{}))
	assert.Nil(t, exitCmd.ParseOptions([]string{"test"}))

}

func TestExecute(t *testing.T) {
	// TODO: Test os.exit() for 100 % coverage

}

func TestPrintUsage(t *testing.T) {
	// should be equal
	var exitCmd *exit.Exit
	assert.Equal(t, exitCmd.PrintUsage(), "Usage: exit ")

}
