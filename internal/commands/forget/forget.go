package forget

import (
	"errors"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"

	"github.com/rs/zerolog/log"
)

type Forget struct {
	hash kademliaid.KademliaID
}

func (forget *Forget) Execute(node *node.Node) (string, error) {
	log.Trace().Msg("Executing forget command")

	log.Trace().Str("ID", forget.hash.String()).Msg("Forgetting data entry")
	err := node.DataStore.Forget(&forget.hash)

	return "", err
}

func (forget *Forget) ParseOptions(options []string) error {
	if len(options) < 1 {
		return errors.New("Missing hash")
	}
	forget.hash = *kademliaid.FromString(options[0])
	return nil
}

func (forget *Forget) PrintUsage() string {
	return "USAGE: forget <hash>"
}
