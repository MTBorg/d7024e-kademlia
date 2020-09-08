package message

import (
	"errors"
	kademliaMessage "kademlia/internal/rpc"
	"kademlia/internal/udpsender"

	"github.com/rs/zerolog/log"
)

type Message struct {
	Target  string
	Content string
}

func (msg Message) Execute() (string, error) {
	log.Debug().Str("Target", msg.Target).Msg("Executing message command")

	message := kademliaMessage.New(msg.Content, msg.Target)
	udpSender := udpsender.New(msg.Target)
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

func (msg *Message) PrintUsage() {
	log.Info().Msg("Usage: msg {target address} {message content}")
}
