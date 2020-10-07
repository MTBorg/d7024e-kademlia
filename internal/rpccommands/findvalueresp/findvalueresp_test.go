package findvalueresp_test

import (
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpccommands/findvalueresp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	senderId := kademliaid.NewRandomKademliaID()
	rpcID := kademliaid.NewRandomKademliaID()

	// should return a new FIND_VALUE_RESP RPC
	findResp := findvalueresp.New(senderId, rpcID)
	assert.IsType(t, findvalueresp.FindValueResp{}, *findResp)
}

func TestParseOptions(t *testing.T) {
	findRespCmd := findvalueresp.FindValueResp{}

	// should report an error if the response contains no data
	res := findRespCmd.ParseOptions(&[]string{})
	assert.Error(t, res)

	// should not report an error if the reponse contains some data
	res = findRespCmd.ParseOptions(&[]string{"blabla"})
	assert.NoError(t, res)

}
