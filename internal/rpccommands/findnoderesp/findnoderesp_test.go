package findenoderesp_test

import (
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpccommands/findnoderesp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var fnResp *findenoderesp.FindNodeResp
	rpcId := kademliaid.NewRandomKademliaID()

	// should return a FindNode object
	fnResp = findenoderesp.New(rpcId)
	assert.IsType(t, &findenoderesp.FindNodeResp{}, fnResp)
}

func TestParseOptions(t *testing.T) {
	var fnResp *findenoderesp.FindNodeResp

	rpcId := kademliaid.NewRandomKademliaID()
	fnResp = findenoderesp.New(rpcId)

	// should report an error if the response contains no data
	options := []string{}
	assert.Error(t, fnResp.ParseOptions(&options))

	// should not report an error if the response contains data
	options = []string{"mydata"}
	assert.NoError(t, fnResp.ParseOptions(&options))
}
