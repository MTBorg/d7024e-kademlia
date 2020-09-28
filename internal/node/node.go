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
	"regexp"
	"strconv"
	"strings"
	"sync"

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
		var contactsMutex sync.Mutex
		for i := 0; i < sl.Len() && numProbed < alpha; i++ {
			if !sl.Entries[i].Probed {
				sl.Entries[i].Probed = true
				rpc := node.NewRPC(fmt.Sprintf("%s %s", "FIND_NODE", id), sl.Entries[i].Contact.Address)
				node.RPCPool.Lock()
				node.RPCPool.Add(rpc.RPCId)
				entryRPC := node.RPCPool.GetEntry(rpc.RPCId)
				node.RPCPool.Unlock()
				rpcIds = append(rpcIds, rpc.RPCId)
				channels[numProbed] = entryRPC.Channel
				numProbed++
				network.Net.SendFindContactMessage(&rpc)
			}
		}

		// If no new nodes were probed this iteration the search is done
		if numProbed == 0 {
			log.Trace().Msg("FIND_NODE lookup became stale")
			break
		}

		// Handle response from probed nodes
		contacts := []*contact.Contact{}
		var wg sync.WaitGroup
		wg.Add(numProbed)
		for i := 0; i < numProbed; i++ {
			go func(i int, wg *sync.WaitGroup, contactsMutex *sync.Mutex) {
				defer wg.Done()
				data := <-channels[i]
				node.RPCPool.Lock()
				node.RPCPool.Delete(rpcIds[i]) // remove from rpc pool
				node.RPCPool.Unlock()

				// parse contacts from response data
				strContacts := strings.Split(data, " ")
				for _, strContact := range strContacts {
					err, contact := contact.Deserialize(&strContact)
					if err == nil {
						contactsMutex.Lock()
						contacts = append(contacts, contact)
						contactsMutex.Unlock()
					}
				}
			}(i, &wg, &contactsMutex)
		}

		wg.Wait()

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

	return sl.GetContacts()
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

func (node *Node) LookupData(hash *kademliaid.KademliaID) string {
	alpha, err := strconv.Atoi(os.Getenv("ALPHA"))
	if err != nil {
		log.Error().Msgf("Failed to convert env variable ALPHA from string to int: %s", err)
	}

	sl := shortlist.NewShortlist(hash, node.FindKClosest(hash, nil, alpha))

	// iterative lookup until the search becomes stale and no closer node
	// can be found
	stale := false
	for !stale {
		channels := make([]chan string, alpha)
		var contactsMutex sync.Mutex

		// probe at most the alpha closest nodes
		probed := 0
		rpcIDs := []*kademliaid.KademliaID{}
		for i := 0; i < sl.Len() && probed < alpha; i++ {
			if !sl.Entries[i].Probed && !sl.Entries[i].Contact.ID.Equals(node.ID) {
				sl.Entries[i].Probed = true
				rpc := node.NewRPC(
					fmt.Sprintf("FIND_VALUE %s", hash.String()),
					sl.Entries[i].Contact.Address)
				rpcIDs = append(rpcIDs, rpc.RPCId)

				node.RPCPool.Lock()
				node.RPCPool.Add(rpc.RPCId)
				entry := node.RPCPool.GetEntry(rpc.RPCId)
				node.RPCPool.Unlock()

				channels[probed] = entry.Channel
				probed++
				network.Net.SendFindDataMessage(&rpc)
			}
		}

		if probed == 0 {
			log.Trace().Msg("FIND_VALUE lookup became stale")
			break
		}

		contacts := []*contact.Contact{}
		result := ""
		var wg sync.WaitGroup
		wg.Add(probed)
		for i := 0; i < probed; i++ {
			go func(i int, wg *sync.WaitGroup, contactsMutex *sync.Mutex) {
				defer wg.Done()

				data := <-channels[i]

				if match, _ := regexp.MatchString("VALUE=.*", data); match { // Value was found
					regex := regexp.MustCompile(`=`)
					s := regex.Split(data, 2)
					value := s[1]
					log.Info().Str("Value", value).Msg("Found value")

					result = value
				} else {
					sContacts := strings.Split(data, " ")
					for _, sContact := range sContacts {
						err, c := contact.Deserialize(&sContact)
						if err == nil {
							c.CalcDistance(hash)
							contactsMutex.Lock()
							contacts = append(contacts, c)
							contactsMutex.Unlock()
						} else {
							log.Warn().Msgf("Failed to deserialize contact: %s", err)
							log.Print(sContact)
						}
					}
				}
			}(i, &wg, &contactsMutex)
		}

		wg.Wait()

		node.RPCPool.Lock()
		for i := 0; i < probed; i++ {
			node.RPCPool.Delete(rpcIDs[i])
		}
		node.RPCPool.Unlock()

		if result != "" {
			return result
		}

		for _, c := range contacts {
			sl.Add(c)
		}
	}

	return ""
}

func (node *Node) Store(value *string) {
	log.Trace().Str("Value", *value).Msg("Storing value")
	node.DataStore.Insert(*value)
}

// FindKClosest returns a list of candidates containing the k closest nodes
// to the key being searched for (from the nodes own bucket(s))
func (node *Node) FindKClosest(key *kademliaid.KademliaID, requestorID *kademliaid.KademliaID, k int) []contact.Contact {
	return node.RoutingTable.FindClosestContacts(key, requestorID, k)
}
