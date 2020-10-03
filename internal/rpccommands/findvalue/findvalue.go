package findvalue

import (
	"errors"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

type FindValue struct {
	hash      *kademliaid.KademliaID
	rpcId     *kademliaid.KademliaID
	requestor *contact.Contact
}

func New(requestor *contact.Contact, rpcId *kademliaid.KademliaID) *FindValue {
	return &FindValue{requestor: requestor, rpcId: rpcId}
}

func (find *FindValue) Execute(node *node.Node) {
	log.Trace().Msg("Executing FIND_VALUE RPC")

	if value := node.DataStore.Get(*find.hash); value != "" {
		log.Debug().Str("Value", value).Str("Hash", find.hash.String()).Msg("Found key")
		s := "VALUE=" + value
		node.Network.SendFindDataRespMessage(node.ID, find.requestor.Address, find.rpcId, &s)
	} else {

		k, err := strconv.Atoi(os.Getenv("K"))
		if err != nil {
			log.Error().Msgf("Failed to convert env variable ALPHA from string to int: %s", err)
		}
		log.Debug().Str("Hash", find.hash.String()).Msg("Did not find key")
		closest := node.FindKClosest(find.hash, find.requestor.ID, k)
		data := contact.SerializeContacts(closest)
		node.Network.SendFindDataRespMessage(node.ID, find.requestor.Address, find.rpcId, &data)
	}
}

func (find *FindValue) ParseOptions(options *[]string) error {
	if (len(*options)) == 0 {
		return errors.New("Missing hash")
	}

	find.hash = kademliaid.FromString((*options)[0])
	return nil
}
