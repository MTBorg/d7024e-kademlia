package rpcparser_test

import (
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpc"
	rpcparser "kademlia/internal/rpc/parser"
	"kademlia/internal/rpccommand"
	"kademlia/internal/rpccommands/ping"
	"kademlia/internal/rpccommands/pong"
	"kademlia/internal/rpccommands/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRPC(t *testing.T) {
	adr := address.New("127.0.0.1:1776")
	c := contact.NewContact(kademliaid.NewRandomKademliaID(), &adr)
	var r rpc.RPC
	var rpcCmd rpccommand.RPCCommand
	var err error

	//Should be able to parse a PING rpc
	r = rpc.New("PING", "sometarget")
	rpcCmd, err = rpcparser.ParseRPC(&c, &r)
	assert.Nil(t, err)
	assert.IsType(t, ping.Ping{}, rpcCmd)

	//Should be able to parse a PONG rpc
	r = rpc.New("PONG", "sometarget")
	rpcCmd, err = rpcparser.ParseRPC(&c, &r)
	assert.Nil(t, err)
	assert.IsType(t, pong.Pong{}, rpcCmd)

	//Should be able to parse a STORE rpc
	r = rpc.New("STORE", "sometarget")
	rpcCmd, err = rpcparser.ParseRPC(&c, &r)
	assert.Nil(t, err)
	assert.IsType(t, &store.Store{}, rpcCmd)

	//Should not parse an unknown RPC
	r = rpc.New("HELLO", "sometarget")
	rpcCmd, err = rpcparser.ParseRPC(&c, &r)
	assert.EqualError(t, err, "Received unknown RPC HELLO")
	assert.Nil(t, rpcCmd)

	//Should not parse empty string
	r = rpc.New("", "sometarget")
	rpcCmd, err = rpcparser.ParseRPC(&c, &r)
	assert.Error(t, err)
	assert.EqualError(t, err, "Missing RPC name")
	assert.Nil(t, rpcCmd)
}
