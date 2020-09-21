package put

import (
	"errors"
	"github.com/rs/zerolog/log"
	"kademlia/internal/kademliaid"
	"kademlia/internal/network"
	"kademlia/internal/node"
	"strings"
)

type Put struct {
	fileContent string
}

func (put *Put) Execute(node *node.Node) (string, error) {
	log.Debug().Msg("Executing put command")

	key := kademliaid.NewKademliaID(&put.fileContent)
	// closestNodes := node.FindKClosest(&key, nil, k)
	closestNodes := node.LookupContact(&key)

	node.Store(&put.fileContent)

	// Send STORE RPCs
	for _, node := range closestNodes {
		network.Net.SendStoreMessage(node.ID, node.Address, []byte(put.fileContent))
	}

	return "", nil
}

func (put *Put) ParseOptions(options []string) error {
	if len(options) < 1 {
		return errors.New("Missing file content")
	}
	put.fileContent = strings.Join(options[0:], " ")
	return nil
}

func (put *Put) PrintUsage() {
	log.Info().Msg("put <file content>")
}
