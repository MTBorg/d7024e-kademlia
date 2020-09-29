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

// Probes at most alpha nodes from the shortlist with content
func (node *Node) probeAlpha(
	sl *shortlist.Shortlist,
	channels *[]chan string,
	content string,
	alpha int) (int, []*kademliaid.KademliaID) {

	numProbed := 0
	rpcIds := []*kademliaid.KademliaID{}
	for i := 0; i < sl.Len() && numProbed < alpha; i++ {
		if !sl.Entries[i].Probed {
			log.Trace().Str("NodeID", sl.Entries[i].Contact.ID.String()).Msg("Probing node")
			sl.Entries[i].Probed = true
			rpc := node.NewRPC(content, sl.Entries[i].Contact.Address)

			node.RPCPool.Lock()
			node.RPCPool.Add(rpc.RPCId)
			entryRPC := node.RPCPool.GetEntry(rpc.RPCId)
			node.RPCPool.Unlock()

			rpcIds = append(rpcIds, rpc.RPCId)
			(*channels)[numProbed] = entryRPC.Channel
			numProbed++
			network.Net.SendFindContactMessage(&rpc)
		}
	}
	return numProbed, rpcIds
}

func deserializeContacts(data string, targetId *kademliaid.KademliaID) []*contact.Contact {
	contacts := []*contact.Contact{}
	for _, sContact := range strings.Split(data, " ") {
		err, c := contact.Deserialize(&sContact)
		if err == nil {
			c.CalcDistance(targetId)
			contacts = append(contacts, c)
		} else {
			log.Warn().Msgf("Failed to deserialize contact: %s", err)
			log.Print(sContact)
		}
	}
	return contacts
}

// Handles the responses from the probed nodes during a node lookup
func (node *Node) lookupContactHandleResponses(
	sl *shortlist.Shortlist,
	targetId *kademliaid.KademliaID,
	numProbed int,
	channels *[]chan string,
	rpcIds []*kademliaid.KademliaID) {
	// Handle response from probed nodes
	var contactsMutex sync.Mutex
	contacts := []*contact.Contact{}
	var wg sync.WaitGroup
	wg.Add(numProbed)
	for i := 0; i < numProbed; i++ {
		go func(i int, wg *sync.WaitGroup, contactsMutex *sync.Mutex) {
			defer wg.Done()
			data := <-(*channels)[i]

			// parse contacts from response data
			contactsMutex.Lock()
			contacts = append(contacts, deserializeContacts(data, targetId)...)
			contactsMutex.Unlock()
		}(i, &wg, &contactsMutex)
	}
	wg.Wait()

	node.RPCPool.Lock()
	for i := 0; i < numProbed; i++ {
		node.RPCPool.Delete(rpcIds[i])
	}
	node.RPCPool.Unlock()

	for _, contact := range contacts {
		sl.Add(contact)
	}
}

func (node *Node) lookupDataHandleResponses(sl *shortlist.Shortlist,
	targetId *kademliaid.KademliaID,
	numProbed int,
	channels *[]chan string,
	rpcIds []*kademliaid.KademliaID) string {

	contacts := []*contact.Contact{}
	var contactsMutex sync.Mutex
	result := ""
	var wg sync.WaitGroup
	wg.Add(numProbed)
	for i := 0; i < numProbed; i++ {
		go func(i int, wg *sync.WaitGroup, contactsMutex *sync.Mutex) {
			defer wg.Done()
			data := <-(*channels)[i]

			if match, _ := regexp.MatchString("VALUE=.*", data); match { // Value was found
				regex := regexp.MustCompile(`=`)

				// Extract value pairs
				fields := strings.Split(data, "/")
				valueField := fields[0]
				senderIdField := fields[1]

				// Extract values
				value := regex.Split(valueField, 2)[1]
				senderId := regex.Split(senderIdField, 2)[1]

				log.Info().Str("Value", value).Str("SenderID", senderId).Msg("Found value")

				result = value + ", from node: " + senderId
				sl.Entries[i].ReturnedValue = true
			} else {
				contactsMutex.Lock()
				contacts = append(contacts, deserializeContacts(data, targetId)...)
				contactsMutex.Unlock()
			}
		}(i, &wg, &contactsMutex)
	}
	wg.Wait()

	node.RPCPool.Lock()
	for i := 0; i < numProbed; i++ {
		node.RPCPool.Delete(rpcIds[i])
	}
	node.RPCPool.Unlock()

	for _, c := range contacts {
		sl.Add(c)
	}

	return result
}

// LookupContact searches for the contact with the specified key using the node
// lookup algorithm.
//
// TODO: Ignore request after waiting X time
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

		numProbed, rpcIds := node.probeAlpha(sl, &channels, fmt.Sprintf("%s %s", "FIND_NODE", id), alpha)

		// If no new nodes were probed this iteration the search is done
		if numProbed == 0 {
			log.Trace().Msg("FIND_NODE lookup became stale")
			break
		}

		node.lookupContactHandleResponses(sl, id, numProbed, &channels, rpcIds)

		// Send FIND_NODE to all unqueried nodes in the shortlist and terminate
		// the search since no node closer to the target was found this iteration
		if sl.Closest == closestSoFar {
			log.Trace().Msg("Closest node not updated")
			numProbed, rpcIds := node.probeAlpha(sl, &channels, fmt.Sprintf("%s %s", "FIND_NODE", id), alpha)

			node.lookupContactHandleResponses(sl, id, numProbed, &channels, rpcIds)
			break
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
	result := ""
	for {
		channels := make([]chan string, alpha)
		closestSoFar := sl.Closest

		numProbed, rpcIDs := node.probeAlpha(sl, &channels, fmt.Sprintf("FIND_VALUE %s", hash.String()), alpha)

		if numProbed == 0 {
			log.Trace().Msg("FIND_VALUE lookup became stale")
			break
		}

		result = node.lookupDataHandleResponses(sl, hash, numProbed, &channels, rpcIDs)
		if result != "" {
			return result
		}

		if sl.Closest == closestSoFar {
			log.Trace().Msg("Closest node not updated")
			numProbed, rpcIds := node.probeAlpha(sl, &channels, fmt.Sprintf("%s %s", "FIND_VALUE", hash), alpha)

			result = node.lookupDataHandleResponses(sl, hash, numProbed, &channels, rpcIds)
			if result != "" {
				return result
			}
		}

	}

	s := "Value not found, k closest contacts: ["
	for i, entry := range sl.Entries {
		s += entry.Contact.String()
		if i < len(sl.Entries)-1 {
			s += ", "
		}
	}
	s += "]"
	return s
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
