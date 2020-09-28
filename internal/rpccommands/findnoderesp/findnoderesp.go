package findenoderesp

import (
	"errors"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"strings"

	"github.com/rs/zerolog/log"
)

type FindNodeResp struct {
	rpcId *kademliaid.KademliaID
	data  *string
}

func New(rpcId *kademliaid.KademliaID) *FindNodeResp {
	return &FindNodeResp{rpcId: rpcId}
}

func (fnResp *FindNodeResp) Execute(node *node.Node) {
	log.Trace().Msg("Executing FIND_NODE_RESPONSE RPC")
	node.RPCPool.Lock()
	entry := node.NodeData.RPCPool.GetEntry(fnResp.rpcId)
	node.RPCPool.Unlock()
	if entry != nil {
		log.Trace().Str("rpcID", fnResp.rpcId.String()).Msg("Writing to channel with rpcID")
		entry.Channel <- *fnResp.data
		log.Trace().Str("rpcID", fnResp.rpcId.String()).Msg("Done writing to channel")
	} else {
		log.Warn().Str("RPCId", fnResp.rpcId.String()).Msg("Received FIND_NODE_RESPONSE with unknown RPCId")
	}
}

func (fnResp *FindNodeResp) ParseOptions(options *[]string) error {
	if len(*options) == 0 {
		return errors.New("Received empty FIND_NODE_RESPONSE RPC")
	}
	data := strings.Join(*options, " ")
	fnResp.data = &data
	return nil
}
