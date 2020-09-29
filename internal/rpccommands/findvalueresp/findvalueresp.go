package findvalueresp

import (
	"errors"
	"kademlia/internal/kademliaid"
	"kademlia/internal/node"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

type FindValueResp struct {
	content  string
	rpcId    *kademliaid.KademliaID
	senderId *kademliaid.KademliaID
}

func New(senderId *kademliaid.KademliaID, rpcId *kademliaid.KademliaID) *FindValueResp {
	return &FindValueResp{senderId: senderId, rpcId: rpcId}
}

func (findresp *FindValueResp) Execute(node *node.Node) {
	log.Trace().Msg("Executing FIND_VALUE_RESP RPC")
	node.RPCPool.Lock()
	entry := node.RPCPool.GetEntry(findresp.rpcId)
	node.RPCPool.Unlock()
	if entry != nil {
		if match, _ := regexp.MatchString("VALUE=.*", findresp.content); match { // Value was found
			entry.Channel <- findresp.content + "/SENDERID=" + findresp.senderId.String()
		} else {
			entry.Channel <- findresp.content
		}
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
