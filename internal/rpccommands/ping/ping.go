package ping

import (
	"kademlia/internal/address"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"

	"github.com/rs/zerolog/log"
)

type Ping struct {
	senderAddress *address.Address
	rpcId         *kademliaid.KademliaID
}

func New(senderAddress *address.Address, rpcId *kademliaid.KademliaID) Ping {
	return Ping{senderAddress: senderAddress, rpcId: rpcId}
}

func (ping Ping) Execute(node *node.Node) {
	log.Trace().Msg("Executing PING RPC")
	// Respond with pong
	node.Network.SendPongMessage(node.ID, ping.senderAddress, ping.rpcId)
}

func (ping Ping) ParseOptions(options *[]string) error {
	// Ping takes no options
	return nil
}
