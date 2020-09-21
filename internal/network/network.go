package network

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpc"
	"kademlia/internal/udpsender"
)

var Net Network

type Network struct{}

// SendPongMessage replies a "PONG" message to the remote "pinger" address
func (network *Network) SendPongMessage(target *address.Address, id *kademliaid.KademliaID) {

	log.Debug().Str("Address", target.String()).Msg("Sending PONG to address")
	rpc := rpc.New("PONG", target)
	rpc.RPCId = id
	udpSender := udpsender.New(target)

	err := rpc.Send(udpSender)
	if err != nil {
		log.Error().Msgf("Failed to write RPC PING message to UDP: %s", err.Error())
		log.Info().Str("Address", target.String()).Str("Content", "PING").Msg("Message sent to address")
	}
}

// SendPingMessage sends a "PING" message to a remote address
func (network *Network) SendPingMessage(target *address.Address) {
	rpc := rpc.New("PING", target)
	udpSender := udpsender.New(target)

	err := rpc.Send(udpSender)
	if err != nil {
		log.Error().Msgf("Failed to write RPC PING message to UDP: %s", err.Error())
		log.Info().Str("Address", target.String()).Str("Content", "PING").Msg("Message sent to address")
	}
}

func (network *Network) SendFindContactMessage(contact *contact.Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(target *address.Address, data []byte) {
	log.Debug().Str("Target", target.String()).Msg("Sending store message")
	rpc := rpc.New(fmt.Sprintf("%s %s", "STORE", data), target)
	udpSender := udpsender.New(target)
	err := rpc.Send(udpSender)

	if err != nil {
		log.Error().Msgf("Failed to write RPC STORE message to UDP: %s", err.Error())
		log.Info().Str("Address", target.String()).Str("Content", "STORE").Msg("Message sent to address")
	}
}
