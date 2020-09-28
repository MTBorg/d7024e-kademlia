package store

import (
	"errors"
	"kademlia/internal/node"
	"strings"

	"github.com/rs/zerolog/log"
)

type Store struct {
	fileContent string
}

func (store *Store) Execute(node *node.Node) {
	log.Trace().Msg("Executing FIND_NODE RPC")
	node.Store(&store.fileContent)
}

func (store *Store) ParseOptions(options *[]string) error {
	if len(*options) == 0 {
		return errors.New("Received empty STORE RPC")
	}
	store.fileContent = strings.Join(*options, " ")
	return nil
}
