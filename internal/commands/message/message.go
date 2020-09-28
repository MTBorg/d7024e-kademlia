package message

import (
	"errors"
	"kademlia/internal/address"
	"kademlia/internal/node"
	kademliaMessage "kademlia/internal/rpc"
	"kademlia/internal/udpsender"

	"github.com/rs/zerolog/log"
)

type Message struct {
	Target  string
	Content string
}

func (msg Message) Execute(node *node.Node) (string, error) {
	log.Trace().Str("Target", msg.Target).Msg("Executing message command")
	adr := address.New(msg.Target)
	message := kademliaMessage.New(node.ID, msg.Content, adr)
	udpSender := udpsender.New(adr)
	err := message.Send(udpSender)

	if err != nil {
		log.Error().Msgf("Failed to write message to UDP: %s", err.Error())
		log.Info().Str("Address", msg.Target).Str("Content", msg.Content).Msg("Message sent to address")
	}

	return "Message sent!", nil
}

func (msg *Message) ParseOptions(options []string) error {
	if len(options) < 2 {
		return errors.New("Missing target address or content in msg command")
	}
	msg.Target = options[0]
	msg.Content = options[1]
	return nil
}

func (msg *Message) PrintUsage() string {
	return "Usage: msg {target address} {message content}"
}
