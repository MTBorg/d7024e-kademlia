package nodedata

import (
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
	"kademlia/internal/routingtable"
	"kademlia/internal/rpcpool"
)

type NodeData struct {
	RoutingTable *routingtable.RoutingTable
	DataStore    datastore.DataStore
	ID           *kademliaid.KademliaID
	RPCPool      *rpcpool.RPCPool
}
