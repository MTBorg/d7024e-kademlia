package findnode

import (
	"errors"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/network"
	"kademlia/internal/node"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

type FindNode struct {
	id        *string
	requestor *contact.Contact
	rpcId     *kademliaid.KademliaID
}

func New(requestor *contact.Contact, rpcId *kademliaid.KademliaID) *FindNode {
	return &FindNode{requestor: requestor, rpcId: rpcId}
}

func (fn *FindNode) Execute(node *node.Node) {
	log.Trace().Msg("Executing FIND_NODE RPC")

	k, err := strconv.Atoi(os.Getenv("K"))
	if err != nil {
		log.Error().Msgf("Failed to convert env variable ALPHA from string to int: %s", err)
	}

	// Respond with k closest nodes to the key
	kClosest := node.FindKClosest(kademliaid.FromString(*fn.id), fn.requestor.ID, k)
	content := contact.SerializeContacts(kClosest)
	network.Net.SendFindContactRespMessage(node.ID, fn.requestor.Address, fn.rpcId, &content)
}

func (fn *FindNode) ParseOptions(options *[]string) error {
	if len(*options) == 0 {
		return errors.New("Recieved empty FIND_NODE RPC, Missing ID argument")
	}
	fn.id = &(*options)[0]
	return nil
}
