package rpcpool

import (
	"kademlia/internal/kademliaid"
	"sync"
)

type Entry struct {
	Channel chan string
	rpcID   *kademliaid.KademliaID
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

func (pool *RPCPool) Add(rpcID *kademliaid.KademliaID) {
	pool.entries[*rpcID] = &Entry{rpcID: rpcID, Channel: make(chan string)}
}

func (pool *RPCPool) GetEntry(rpcId *kademliaid.KademliaID) *Entry {
	return pool.entries[*rpcId]
}

func (pool *RPCPool) Delete(rpcId *kademliaid.KademliaID) {
	delete(pool.entries, *rpcId)
}

func (pool *RPCPool) Lock() {
	pool.lock.Lock()
}

func (pool *RPCPool) Unlock() {
	pool.lock.Unlock()
}
