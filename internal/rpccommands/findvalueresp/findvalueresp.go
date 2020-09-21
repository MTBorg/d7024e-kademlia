package findvalueresp

import (
	"errors"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"strings"

	"github.com/rs/zerolog/log"
)

type FindValueResp struct {
	content string
	rpcId   *kademliaid.KademliaID
}

func New(rpcId *kademliaid.KademliaID) *FindValueResp {
	return &FindValueResp{rpcId: rpcId}
}

func (findresp *FindValueResp) Execute(node *node.Node) {
	log.Trace().Msg("Executing FIND_VALUE_RESP RPC")
	entry := node.RPCPool.GetEntry(findresp.rpcId)
	if entry != nil {
		entry.Channel <- findresp.content
	} else {
		log.Warn().Str("RPCID", findresp.rpcId.String()).Msg("Tried write to nil RPCPool")
	}
}

func (findresp *FindValueResp) ParseOptions(options *[]string) error {
	if len(*options) == 0 {
		return errors.New("Missing content")
	}
	findresp.content = strings.Join((*options)[0:], " ")
	return nil
}
