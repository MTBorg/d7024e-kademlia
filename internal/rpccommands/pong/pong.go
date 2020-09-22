package pong

import "kademlia/internal/node"

type Pong struct{}

func New() Pong {
	return Pong{}
}

func (pong Pong) Execute(node *node.Node) {
	// Pong does nothing
}

func (pong Pong) ParseOptions(options *[]string) error {
	// Pong takes no options
	return nil
}
