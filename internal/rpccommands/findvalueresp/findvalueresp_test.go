package findvalueresp_test

import (
	"kademlia/internal/contact"
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"kademlia/internal/nodedata"
	"kademlia/internal/rpccommands/findvalueresp"
	"kademlia/internal/rpcpool"
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

func TestExecute(t *testing.T) {
	n := node.Node{NodeData: nodedata.NodeData{DataStore: datastore.New(), RPCPool: rpcpool.New()}}
	senderId := kademliaid.NewRandomKademliaID()
	rpcID := kademliaid.NewRandomKademliaID()
	contact := contact.Contact{ID: senderId}

	n.RPCPool.Add(rpcID, &contact)
	findRespCmd := findvalueresp.New(senderId, rpcID)
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

	// Should return the content and senderId if the value was found
	content = "VALUE=some value" + senderId.String()
	findRespCmd.ParseOptions(&[]string{content})
	go func() {
		findRespCmd.Execute(&n)
	}()
	channel = n.RPCPool.GetEntry(rpcID).Channel
	data = <-channel
	assert.Equal(t, content+"/SENDERID="+senderId.String(), data)
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
