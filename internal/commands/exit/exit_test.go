package exit_test

import (
	"kademlia/internal/commands/exit"
	"kademlia/internal/datastore"
	"kademlia/internal/node"
	"kademlia/internal/nodedata"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParseOptions(t *testing.T) {
	// should be nil
	var exitCmd *exit.Exit
	assert.Nil(t, exitCmd.ParseOptions([]string{}))
	assert.Nil(t, exitCmd.ParseOptions([]string{"test"}))

}

type ExitMock struct {
	mock.Mock
}

func (m *ExitMock) exit(n int) {
	m.Called(0)
}

func TestExecute(t *testing.T) {
	exitMock := ExitMock{}
	exitMock.On("exit", 0)
	exit.ExitFunction = exitMock.exit
	e := exit.Exit{}
	n := node.Node{NodeData: nodedata.NodeData{DataStore: datastore.New()}}
	e.Execute(&n)
	exitMock.AssertExpectations(t)
}

func TestPrintUsage(t *testing.T) {
	// should be equal
	var exitCmd *exit.Exit
	assert.Equal(t, exitCmd.PrintUsage(), "Usage: exit ")

}
