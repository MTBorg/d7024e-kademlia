package rpc

import (
	"errors"
	"fmt"
	"strings"

	"kademlia/internal/address"
	"kademlia/internal/globals"
	"kademlia/internal/kademliaid"

	"github.com/rs/zerolog/log"
)

type RPC struct {
	SenderId *kademliaid.KademliaID
	RPCId    *kademliaid.KademliaID
	Content  string
	Target   *address.Address
}

type Sender interface {
	Send(string) error
}

func New(content string, target *address.Address) RPC {
	return RPC{SenderId: globals.ID, RPCId: kademliaid.NewRandomKademliaID(), Content: content, Target: target}
}

// Constructs a new RPC with a given rpcID.
//
// Useful for creating new RPC's that are responses to previous RPCs, and thus
// should use the same RPCId.
func NewWithID(content string, target *address.Address, rpcId *kademliaid.KademliaID) RPC {
	return RPC{
		SenderId: globals.ID,
		RPCId:    rpcId,
		Content:  content,
		Target:   target,
	}
}

// Sends the message using the send function
func (rpc *RPC) Send(sender Sender) error {
	return sender.Send(rpc.serialize())
}

func (rpc *RPC) serialize() string {
	return fmt.Sprintf("%s;%s;%s", rpc.SenderId, rpc.RPCId, rpc.Content)
}

func Deserialize(s string) (RPC, error) {
	log.Debug().Str("String", s).Msg("Dezerializing string")
	fields := strings.Split(s, ";")
	if len(fields) <= 2 {
		return RPC{}, errors.New("Missing sender id or rpc id")
	} else {
		id := kademliaid.FromString(fields[0])
		RPCId := kademliaid.FromString(fields[1])
		return RPC{SenderId: id, RPCId: RPCId, Content: fields[2]}, nil
	}
}
