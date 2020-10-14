package refresh

import (
	"errors"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"

	"github.com/rs/zerolog/log"
)

type Refresh struct {
	hash kademliaid.KademliaID
}

func (r *Refresh) Execute(node *node.Node) {
	log.Trace().Msg("Executing REFRESH RPC")

	// Calling Get restarts the timer
	if value := node.DataStore.Get(r.hash); value == "" {
		log.Warn().Str("Hash", r.hash.String()).Msg("Received refresh on non existent value")
	}
}

func (r *Refresh) ParseOptions(options *[]string) error {
	if len(*options) < 1 {
		return errors.New("Missing hash")
	}
	hash := (*options)[0]
	r.hash = *kademliaid.FromString(hash)
	return nil
}
