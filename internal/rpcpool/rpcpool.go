package rpcpool

import (
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"sync"
)

type Entry struct {
	Channel chan string
	rpcID   *kademliaid.KademliaID
	Contact *contact.Contact
}

type RPCPool struct {
	lock    sync.Mutex
	entries map[kademliaid.KademliaID]*Entry
}

func New() *RPCPool {
	return &RPCPool{
		entries: make(map[kademliaid.KademliaID]*Entry),
	}
}

func (pool *RPCPool) Add(rpcID *kademliaid.KademliaID, contact *contact.Contact) {
	pool.entries[*rpcID] = &Entry{rpcID: rpcID, Contact: contact, Channel: make(chan string)}
}

func (pool *RPCPool) GetEntry(rpcId *kademliaid.KademliaID) *Entry {
	return pool.entries[*rpcId]
}

func (pool *RPCPool) Delete(rpcId *kademliaid.KademliaID) {
	delete(pool.entries, *rpcId)
}

func (pool *RPCPool) WithLock(f func()) {
	pool.lock.Lock()
	f()
	pool.lock.Unlock()
}
