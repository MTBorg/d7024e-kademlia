package network

import (
	"fmt"

	"kademlia/internal/address"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpc"
	"kademlia/internal/udpsender"

	"github.com/rs/zerolog/log"
)

var Net Network

type Network struct{}

// SendPongMessage replies a "PONG" message to the remote "pinger" address
func (network *Network) SendPongMessage(senderId *kademliaid.KademliaID, target *address.Address, id *kademliaid.KademliaID) {

	log.Debug().Str("Address", target.String()).Msg("Sending PONG to address")
	rpc := rpc.New(senderId, "PONG", target)
	rpc.RPCId = id
	udpSender := udpsender.New(target)

	err := rpc.Send(udpSender)
	if err != nil {
		log.Error().Msgf("Failed to write RPC PING message to UDP: %s", err.Error())
		log.Info().Str("Address", target.String()).Str("Content", "PING").Msg("Message sent to address")
	}
}

// SendPingMessage sends a "PING" message to a remote address
func (network *Network) SendPingMessage(senderId *kademliaid.KademliaID, target *address.Address) {
	rpc := rpc.New(senderId, "PING", target)
	udpSender := udpsender.New(target)

	err := rpc.Send(udpSender)
	if err != nil {
		log.Error().Msgf("Failed to write RPC PING message to UDP: %s", err.Error())
		log.Info().Str("Address", target.String()).Str("Content", "PING").Msg("Message sent to address")
	}
}

func (network *Network) SendFindContactMessage(rpc *rpc.RPC) {
	udpSender := udpsender.New(rpc.Target)
	err := rpc.Send(udpSender)
	if err != nil {
		log.Error().Msgf("Failed to write FIND_NODE RPC to UDP: %s", err.Error())
	}
	log.Info().Str("Address", rpc.Target.String()).Str("rpcId", rpc.RPCId.String()).Str("Content", rpc.Content).Msg("FIND_NODE sent to address")
}

// SendFindContactRespMessage responds to a FIND_NODE RPC by returning the k
// closest contacts to the key that the node knows of
func (network *Network) SendFindContactRespMessage(senderId *kademliaid.KademliaID, target *address.Address, rpcId *kademliaid.KademliaID, content *string) {

	rpc := rpc.NewWithID(senderId, fmt.Sprintf("%s %s", "FIND_NODE_RESPONSE", *content), target, rpcId)

	udpSender := udpsender.New(target)
	err := rpc.Send(udpSender)
	if err != nil {
		log.Error().Msgf("Failed to write FIND_NODE_RESPONSE message to UDP: %s", err.Error())
	}
	log.Info().Str("Address", target.String()).Str("Content", *content).Msg("FIND_NODE_RESPONSE sent to address")
}

func (network *Network) SendFindDataMessage(rpc rpc.RPC) {
	//TODO
}

func (network *Network) SendFindDataRespMessage(target *address.Address, rpcId *kademliaid.KademliaID) {
	//TODO
}

func (network *Network) SendStoreMessage(senderId *kademliaid.KademliaID, target *address.Address, data []byte) {
	log.Debug().Str("Target", target.String()).Msg("Sending store message")
	rpc := rpc.New(senderId, fmt.Sprintf("%s %s", "STORE", data), target)
	udpSender := udpsender.New(target)
	err := rpc.Send(udpSender)

	if err != nil {
		log.Error().Msgf("Failed to write RPC STORE message to UDP: %s", err.Error())
		log.Info().Str("Address", target.String()).Str("Content", "STORE").Msg("Message sent to address")
	}
}
