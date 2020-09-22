package rpccommand

import (
	"kademlia/internal/node"
)

type RPCCommand interface {
	Execute(node *node.Node)
	ParseOptions(options *[]string) error
}
