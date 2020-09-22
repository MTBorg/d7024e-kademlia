package ping

import (
	"kademlia/internal/address"
	"kademlia/internal/kademliaid"
	"kademlia/internal/network"
	"kademlia/internal/node"
)

type Ping struct {
	senderAddress *address.Address
	rpcId         *kademliaid.KademliaID
}

func New(senderAddress *address.Address, rpcId *kademliaid.KademliaID) Ping {
	return Ping{senderAddress: senderAddress, rpcId: rpcId}
}

func (ping Ping) Execute(node *node.Node) {
	// Respond with pong
	network.Net.SendPongMessage(node.ID, ping.senderAddress, ping.rpcId)
}

func (ping Ping) ParseOptions(options *[]string) error {
	// Ping takes no options
	return nil
}
