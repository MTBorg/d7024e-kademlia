package rpcparser

import (
	"errors"
	"fmt"
	"kademlia/internal/contact"
	"kademlia/internal/rpc"
	"kademlia/internal/rpccommand"
	"kademlia/internal/rpccommands/findnode"
	"kademlia/internal/rpccommands/findnoderesp"
	"kademlia/internal/rpccommands/findvalue"
	"kademlia/internal/rpccommands/findvalueresp"
	"kademlia/internal/rpccommands/ping"
	"kademlia/internal/rpccommands/pong"
	"kademlia/internal/rpccommands/store"
	"strings"

	"github.com/rs/zerolog/log"
)

// Parses a rpc and returns a rpc command.
func ParseRPC(requestor *contact.Contact, rpc *rpc.RPC) (rpccommand.RPCCommand, error) {
	fields := strings.Fields(rpc.Content)
	if len(fields) == 0 {
		return nil, errors.New("Missing RPC name")
	}

	var cmd rpccommand.RPCCommand
	var err error
	rpcLog := log.Debug().Str("RPCId", rpc.RPCId.String())
	switch identifier := fields[0]; identifier {
	case "PING":
		rpcLog.Msg("PING received")
		cmd = ping.New(requestor.Address, rpc.RPCId)
	case "PONG":
		rpcLog.Msg("PONG received")
		cmd = pong.New()
	case "STORE":
		rpcLog.Msg("STORE received")
		cmd = new(store.Store)
	case "FIND_NODE":
		rpcLog.Msg("FIND_NODE received")
		cmd = findnode.New(requestor, rpc.RPCId)
	case "FIND_NODE_RESPONSE":
		rpcLog.Msg("FIND_NODE_RESPONSE received")
		cmd = findenoderesp.New(rpc.RPCId)
	case "FIND_VALUE":
		rpcLog.Msg("FIND_VALUE received")
		cmd = findvalue.New(requestor, rpc.RPCId)
	case "FIND_VALUE_RESP":
		rpcLog.Msg("FIND_VALUE_RESP received")
		cmd = findvalueresp.New(rpc.RPCId)
	default:
		err = errors.New(fmt.Sprintf("Received unknown RPC %s", identifier))
		cmd = nil
	}
	return cmd, err
}
