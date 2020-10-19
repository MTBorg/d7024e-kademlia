package findenoderesp_test

import (
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"kademlia/internal/nodedata"
	"kademlia/internal/rpccommands/findnoderesp"
	"kademlia/internal/rpcpool"
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

func TestExecute(t *testing.T) {
	n := node.Node{NodeData: nodedata.NodeData{DataStore: datastore.New(), RPCPool: rpcpool.New()}}
	rpcID := kademliaid.NewRandomKademliaID()
	n.RPCPool.Add(rpcID)
	findRespCmd := findenoderesp.New(rpcID)
	var content, data string
	var channel chan string

	// Should return just the content of the rpc if the value was not found
	content = "some content"
	findRespCmd.ParseOptions(&[]string{content})
	go func() {
		findRespCmd.Execute(&n)
	}()
	channel = n.RPCPool.GetEntry(rpcID).Channel
	data = <-channel
	assert.Equal(t, content, data)
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
