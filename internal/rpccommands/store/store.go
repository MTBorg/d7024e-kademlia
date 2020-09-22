package store

import (
	"errors"
	"kademlia/internal/node"
	"strings"
)

type Store struct {
	fileContent string
}

func (store *Store) Execute(node *node.Node) {
	node.Store(&store.fileContent)
}

func (store *Store) ParseOptions(options *[]string) error {
	if len(*options) == 0 {
		return errors.New("Received empty STORE RPC")
	}
	store.fileContent = strings.Join(*options, " ")
	return nil
}
