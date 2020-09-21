package globals

import (
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpcpool"
)

var ID *kademliaid.KademliaID = kademliaid.NewRandomKademliaID()
var RPCPool *rpcpool.RPCPool = rpcpool.New()
