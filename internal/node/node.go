package node

import (
	. "kademlia/internal/contact"
	"kademlia/internal/kademliaid"
)

type Node struct {
	Id *kademliaid.KademliaID
}

var KadNode Node

func (node *Node) Init() {
	KadNode = Node{Id: kademliaid.NewRandomKademliaID()}
}

func (node *Node) LookupContact(target *Contact) {
	// TODO
}

func (node *Node) LookupData(hash string) {
	// TODO
}

func (node *Node) Store(data []byte) {
	// TODO
}
