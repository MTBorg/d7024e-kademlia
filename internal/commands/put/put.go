package put

import (
	"errors"
	"kademlia/internal/kademliaid"
	"kademlia/internal/network"
	"kademlia/internal/node"
	"strings"

	"github.com/rs/zerolog/log"
)

type Put struct {
	fileContent string
}

func (put *Put) Execute() (string, error) {
	log.Debug().Msg("Executing put command")
	k := 20 //TODO: Use constant
	key := kademliaid.NewKademliaID(&put.fileContent)
	closestNodes := node.KadNode.RoutingTable.FindClosestContacts(&key, k)

	node.KadNode.Store(&put.fileContent)

	// Send STORE RPCs
	for _, node := range closestNodes {
		network.Net.SendStoreMessage(node.Address, []byte(put.fileContent))
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
