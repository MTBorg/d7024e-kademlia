package node

import (
	"fmt"
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
	"kademlia/internal/network"
	"kademlia/internal/nodedata"
	"kademlia/internal/routingtable"
	"kademlia/internal/rpc"
	"kademlia/internal/rpcpool"
	"kademlia/internal/shortlist"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type Node struct {
	nodedata.NodeData
}

// Initialize the node by generating a NodeID and creating a new routing table
// containing itself as a contact
func (node *Node) Init(address *address.Address) {
	id := kademliaid.NewRandomKademliaID()
	me := contact.NewContact(id, address)
	*node = Node{
		NodeData: nodedata.NodeData{
			RoutingTable: routingtable.NewRoutingTable(me),
			DataStore:    datastore.New(),
			ID:           id,
			RPCPool:      rpcpool.New(),
		},
	}
}

// Join performs a node lookup on itself to join the network and fill its
// routing table
// TODO: Should maybe check that the node knows of another node in the network.
// This is not a problem as long as the init script is used.
func (node *Node) JoinNetwork() {
	node.LookupContact(node.RoutingTable.GetMe().ID)
}

// LookupContact searches for the contact with the specified key using the node
// lookup algorithm.
//
// This implementation uses waitGroups to send alpha parallel FIND_NODE RPCs
// and waits for the response of each request.
// TODO: Ignore request after waiting X time and continue with next iteration
func (node *Node) LookupContact(id *kademliaid.KademliaID) []contact.Contact {
	alpha, err := strconv.Atoi(os.Getenv("ALPHA"))
	if err != nil {
		log.Error().Msgf("Failed to convert env variable ALPHA from string to int: %s", err)
	}

	sl := shortlist.NewShortlist(id, node.FindKClosest(id, nil, alpha))
	channels := make([]chan string, alpha)

	// iterative lookup until the search becomes stale
	for {
		closestSoFar := sl.Closest

		// Probe alpha new nodes
		numProbed := 0
		rpcIds := []*kademliaid.KademliaID{}
		for i := 0; i < sl.Len() && numProbed < alpha; i++ {
			if !sl.Entries[i].Probed {
				sl.Entries[i].Probed = true
				rpc := node.NewRPC(fmt.Sprintf("%s %s", "FIND_NODE", id), sl.Entries[i].Contact.Address)
				node.RPCPool.Add(rpc.RPCId)
				entryRPC := node.RPCPool.GetEntry(rpc.RPCId)
				rpcIds = append(rpcIds, rpc.RPCId)
				channels[numProbed] = entryRPC.Channel
				numProbed++
				log.Trace().Str("Entry", fmt.Sprintf("%d", i)).Msg("Probing entry")
				network.Net.SendFindContactMessage(&rpc)
			}
		}

		// If no new nodes were probed this iteration the search is done
		if numProbed == 0 {
			log.Debug().Msg("Lookup Node became stale")
			break
		}

		// Handle response from probed nodes
		contacts := []*contact.Contact{}
		for i := 0; i < numProbed; i++ {
			log.Trace().Str("Channel", fmt.Sprintf("%d", i)).Str("rpcID", rpcIds[i].String()).Msg("Waiting for channel")
			data := <-channels[i]
			node.RPCPool.Delete(rpcIds[i]) // remove from rpc pool
			log.Trace().Str("Channel", fmt.Sprintf("%d", i)).Str("Data", data).Msg("Received data from channel")

			// parse contacts from response data
			strContacts := strings.Split(data, " ")
			for _, strContact := range strContacts {
				err, contact := contact.Deserialize(&strContact)
				if err != nil {
					log.Print(err)
					log.Warn().Msg("Received FIND_NODE_RESPONSE no contacts in FIND_NODE_RESPONSE")
				} else {
					log.Debug().Str("Contact", contact.String()).Msg("Received a contact")
					contacts = append(contacts, contact)
				}
			}
		}

		// Update shortlist with new contacts recieved from responses
		for _, contact := range contacts {
			sl.Add(contact)
		}

		// Send FIND_NODE to all unqueried nodes in the shortlist and terminate
		// the search since no node closer to the target was found this iteration
		if sl.Closest == closestSoFar {
			// TODO
		}
	}

	log.Debug().Msg("Lookup Node became stale, returning k closest found")
	log.Debug().Msg("k closest contacts found: ")
	contacts := sl.GetContacts()
	s := ""
	for _, c := range contacts {
		s += c.String() + "\n"
	}
	log.Debug().Msg(s)

	return contacts
}

func (node *Node) NewRPC(content string, target *address.Address) rpc.RPC {
	return rpc.RPC{SenderId: node.ID, RPCId: kademliaid.NewRandomKademliaID(), Content: content, Target: target}
}

// Constructs a new RPC with a given rpcID.
//
// Useful for creating new RPC's that are responses to previous RPCs, and thus
// should use the same RPCId.
func NewRPCWithID(senderId *kademliaid.KademliaID, content string, target *address.Address, rpcId *kademliaid.KademliaID) rpc.RPC {
	return rpc.RPC{
		SenderId: senderId,
		RPCId:    rpcId,
		Content:  content,
		Target:   target,
	}
}

func (node *Node) LookupData(hash string) {
	// TODO
}

func (node *Node) Store(value *string) {
	log.Debug().Str("Value", *value).Msg("Storing value")
	node.DataStore.Insert(*value)
}

// FindKClosest returns a list of candidates containing the k closest nodes
// to the key being searched for (from the nodes own bucket(s))
func (node *Node) FindKClosest(key *kademliaid.KademliaID, requestorID *kademliaid.KademliaID, k int) []contact.Contact {
	return node.RoutingTable.FindClosestContacts(key, requestorID, k)
}
