package rpcparser_test

import (
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpc"
	"kademlia/internal/rpc/parser"
	"kademlia/internal/rpccommand"
	"kademlia/internal/rpccommands/findnode"
	"kademlia/internal/rpccommands/findnoderesp"
	"kademlia/internal/rpccommands/findvalue"
	"kademlia/internal/rpccommands/findvalueresp"
	"kademlia/internal/rpccommands/ping"
	"kademlia/internal/rpccommands/pong"
	"kademlia/internal/rpccommands/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRPC(t *testing.T) {
	adr := address.New("127.0.0.1:1776")
	c := contact.NewContact(kademliaid.NewRandomKademliaID(), adr)
	senderId := kademliaid.NewRandomKademliaID()
	var r rpc.RPC
	var rpcCmd rpccommand.RPCCommand
	var err error

	//Should be able to parse a PING rpc
	r = rpc.New(senderId, "PING", adr)
	rpcCmd, err = rpcparser.ParseRPC(&c, &r)
	assert.Nil(t, err)
	assert.IsType(t, ping.Ping{}, rpcCmd)

	//Should be able to parse a PONG rpc
	r = rpc.New(senderId, "PONG", adr)
	rpcCmd, err = rpcparser.ParseRPC(&c, &r)
	assert.Nil(t, err)
	assert.IsType(t, pong.Pong{}, rpcCmd)

	//Should be able to parse a STORE rpc
	r = rpc.New(senderId, "STORE", adr)
	rpcCmd, err = rpcparser.ParseRPC(&c, &r)
	assert.Nil(t, err)
	assert.IsType(t, &store.Store{}, rpcCmd)

	// Should be able to parse a FIND_NODE RPC
	r = rpc.New(senderId, "FIND_NODE", adr)
	rpcCmd, err = rpcparser.ParseRPC(&c, &r)
	assert.Nil(t, err)
	assert.IsType(t, &findnode.FindNode{}, rpcCmd)

	// Should be able to parse a FIND_NODE_RESPONSE RPC
	r = rpc.New(senderId, "FIND_NODE_RESPONSE", adr)
	rpcCmd, err = rpcparser.ParseRPC(&c, &r)
	assert.Nil(t, err)
	assert.IsType(t, &findenoderesp.FindNodeResp{}, rpcCmd)

	// Should be able to parse a FIND_VALUE RPC
	r = rpc.New(senderId, "FIND_VALUE", adr)
	rpcCmd, err = rpcparser.ParseRPC(&c, &r)
	assert.Nil(t, err)
	assert.IsType(t, &findvalue.FindValue{}, rpcCmd)

	// Should be able to parse a FIND_VALUE_RESPONSE RPC
	r = rpc.New(senderId, "FIND_VALUE_RESP", adr)
	rpcCmd, err = rpcparser.ParseRPC(&c, &r)
	assert.Nil(t, err)
	assert.IsType(t, &findvalueresp.FindValueResp{}, rpcCmd)

	//Should not parse an unknown RPC
	r = rpc.New(senderId, "HELLO", adr)
	rpcCmd, err = rpcparser.ParseRPC(&c, &r)
	assert.EqualError(t, err, "Received unknown RPC HELLO")
	assert.Nil(t, rpcCmd)

	//Should not parse empty string
	r = rpc.New(senderId, "", adr)
	rpcCmd, err = rpcparser.ParseRPC(&c, &r)
	assert.Error(t, err)
	assert.EqualError(t, err, "Missing RPC name")
	assert.Nil(t, rpcCmd)
}
