package message

import (
	"errors"
	"fmt"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"strings"
)

type Message struct {
	SenderId *kademliaid.KademliaID
	Content  string
	Target   string
}

type Sender interface {
	Send(string) error
}

func New(content string, target string) Message {
	return Message{SenderId: node.KadNode.Id, Content: content, Target: target}
}

// Sends the message using the send function
func (msg *Message) Send(sender Sender) error {
	return sender.Send(msg.serialize())
}

func (msg *Message) serialize() string {
	return fmt.Sprintf("%s;%s", msg.SenderId, msg.Content)
}

func Deserialize(s string) (Message, error) {
	fields := strings.Split(s, ";")
	if len(fields) <= 1 {
		return Message{}, errors.New("Missing sender id")
	} else {
		id := kademliaid.FromString(fields[0])
		return Message{SenderId: id, Content: fields[1]}, nil
	}
}
