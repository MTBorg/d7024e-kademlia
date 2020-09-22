package node

import (
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
	"kademlia/internal/nodedata"
	"kademlia/internal/routingtable"
	"kademlia/internal/rpc"

	"github.com/rs/zerolog/log"
)

type Node struct {
	nodedata.NodeData
}

// Initialize the node by generating a NodeID and creating a new routing table
// containing itself as a contact
func (node *Node) Init(address address.Address) {
	id := kademliaid.NewRandomKademliaID()
	me := contact.NewContact(id, &address)
	*node = Node{
		NodeData: nodedata.NodeData{
			RoutingTable: routingtable.NewRoutingTable(me),
			DataStore:    datastore.New(),
			ID:           id,
		},
	}
}

func (node *Node) LookupContact(target *contact.Contact) {
	// TODO
}

func (node *Node) NewRPC(content string, target *address.Address) rpc.RPC {
	return rpc.RPC{SenderId: node.ID, RPCId: kademliaid.NewRandomKademliaID(), Content: content, Target: target}
}

// Constructs a new RPC with a given rpcID.
//
// Useful for creating new RPC's that are responses to previous RPCs, and thus
// should use the same RPCId.
func NewRPCWithID(senderId *kademliaid.KademliaID, content string, target *address.Address, rpcId *kademliaid.KademliaID) rpc.RPC {
	return rpc.RPC{
		SenderId: senderId,
		RPCId:    rpcId,
		Content:  content,
		Target:   target,
	}
}

func (node *Node) LookupData(hash string) {
	// TODO
}

func (node *Node) Store(value *string) {
	log.Debug().Str("Value", *value).Msg("Storing value")
	node.DataStore.Insert(*value)
}
