package node

import (
	"fmt"
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/datastore"
	"kademlia/internal/kademliaid"
	"kademlia/internal/network"
	"kademlia/internal/nodedata"
	"kademlia/internal/refreshtimer"
	"kademlia/internal/routingtable"
	"kademlia/internal/rpc"
	"kademlia/internal/rpcpool"
	"kademlia/internal/shortlist"
	"kademlia/internal/udpsender"
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

// Init initializes the node by generating a NodeID and creating a routing
// table, data store and a RPC pool
func (node *Node) Init(address *address.Address) {
	id := kademliaid.NewRandomKademliaID()
	me := contact.NewContact(id, address)
	refreshTimers := []*refreshtimer.RefreshTimer{}
	udpSender, err := udpsender.New()
	if err != nil {
		log.Fatal().Str("Error", err.Error()).Msg("Failed to initialize ndoe")
	}
	*node = Node{
		NodeData: nodedata.NodeData{
			RoutingTable:  routingtable.NewRoutingTable(me),
			DataStore:     datastore.New(),
			ID:            id,
			RPCPool:       rpcpool.New(),
			RefreshTimers: refreshTimers,
			Network:       network.Network{UdpSender: udpSender},
		},
	}

	// start new refresh timers for each bucket, skip last bucket since no node
	// will be in it whp
	for i := 0; i < kademliaid.IDLength*8-1; i++ {
		rt := refreshtimer.NewRefreshTimer(i)
		node.RefreshTimers = append(node.RefreshTimers, rt)
		rt.StartRefreshTimer(node.RefreshBucket)
	}
}

// JoinNetwork performs a node lookup on on the nodes own ID. It then refreshes
// all of its buckets further away than the bucket its closest neighbour is in.
//
// TODO: Should maybe check that the node knows of another node in the network.
// This is not a problem as long as the init script is used.
func (node *Node) JoinNetwork() {

	// lookup on self
	kClosest := node.LookupContact(node.RoutingTable.GetMe().ID)
	if len(kClosest) == 0 {
		log.Error().Msg("Failed to join network: Lookup on self resulted in no contacts")
		return
	}

	// Refresh all buckets further away than the closest neighbour. The first
	// contact will be in the closest neihbour bucket since the list is sorted.
	// Note that the contacts in bucket 0 is the furthest away
	CNBucket := node.RoutingTable.GetBucketIndex(kClosest[0].ID)
	log.Trace().Str("CNBucket", fmt.Sprint(CNBucket)).Msg("Found CNBucket")
	for i := CNBucket - 1; i >= 0; i-- {
		node.RefreshBucket(i)
	}
}

// RefreshBucket refreshes bucket nr bucketIndex by performing a LookupContact
// on a random ID inside the range of the bucket
func (node *Node) RefreshBucket(bucketIndex int) {
	id := kademliaid.NewKademliaIDInRange(node.ID, bucketIndex)
	log.Trace().Str("Bucket", fmt.Sprint(bucketIndex)).Str("IDInRange", id.String()).Msg("Refreshing bucket")
	node.LookupContact(id)
}

// Probes at most alpha nodes from the shortlist with content
func (node *Node) ProbeAlpha(
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

			var entryRPC *rpcpool.Entry
			node.RPCPool.WithLock(func() {
				node.RPCPool.Add(rpc.RPCId)
				entryRPC = node.RPCPool.GetEntry(rpc.RPCId)
			})

			rpcIds = append(rpcIds, rpc.RPCId)
			(*channels)[numProbed] = entryRPC.Channel
			numProbed++
			node.Network.SendFindContactMessage(&rpc)
		}
	}
	return numProbed, rpcIds
}

func DeserializeContacts(data string, targetId *kademliaid.KademliaID) []*contact.Contact {
	contacts := []*contact.Contact{}
	for _, sContact := range strings.Split(data, " ") {
		if sContact != "" {
			err, c := contact.Deserialize(&sContact)
			if err == nil {
				c.CalcDistance(targetId)
				contacts = append(contacts, c)
			}
		}
	}
	return contacts
}

// Handles the responses from the probed nodes during a node lookup
func (node *Node) LookupContactHandleResponses(
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
			contacts = append(contacts, DeserializeContacts(data, targetId)...)
			contactsMutex.Unlock()
		}(i, &wg, &contactsMutex)
	}
	wg.Wait()

	node.RPCPool.WithLock(func() {
		for i := 0; i < numProbed; i++ {
			node.RPCPool.Delete(rpcIds[i])
		}
	})

	for _, contact := range contacts {
		sl.Add(contact)
	}
}

func (node *Node) LookupDataHandleResponses(sl *shortlist.Shortlist,
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
				contacts = append(contacts, DeserializeContacts(data, targetId)...)
				contactsMutex.Unlock()
			}
		}(i, &wg, &contactsMutex)
	}
	wg.Wait()

	node.RPCPool.WithLock(func() {
		for i := 0; i < numProbed; i++ {
			node.RPCPool.Delete(rpcIds[i])
		}
	})

	for _, c := range contacts {
		sl.Add(c)
	}

	return result
}

func GetEnvIntVariable(variable string, defaultValue int) int {
	val, err := strconv.Atoi(os.Getenv(variable))
	if err != nil {
		log.Error().Msgf("Failed to convert env variable %s from string to int: %s", variable, err)
		return defaultValue
	}
	return val
}

// LookupContact searches for the contact with the specified key using the node
// lookup algorithm.
//
// TODO: Ignore request after waiting X time
func (node *Node) LookupContact(id *kademliaid.KademliaID) []contact.Contact {
	alpha, k, sl, channels := setupLookUpAlgorithm(node, id)

	// Restart refresh timer of the bucket this ID is in range of
	if *id != *node.ID {
		bucketIndex := node.RoutingTable.GetBucketIndex(id)
		node.RefreshTimers[bucketIndex].RestartRefreshTimer()
	}

	// iterative lookup until the search becomes stale
	for {
		closestSoFar := sl.Closest

		numProbed, rpcIds := node.ProbeAlpha(sl, &channels, fmt.Sprintf("%s %s", "FIND_NODE", id), alpha)

		// If no new nodes were probed this iteration the search is done
		if numProbed == 0 {
			log.Trace().Msg("FIND_NODE lookup became stale")
			break
		}

		node.LookupContactHandleResponses(sl, id, numProbed, &channels, rpcIds)

		// Send FIND_NODE to all unqueried nodes in the shortlist and terminate
		// the search since no node closer to the target was found this iteration
		if sl.Closest == closestSoFar {
			log.Trace().Msg("Closest node not updated")
			numProbed, rpcIds := node.ProbeAlpha(sl, &channels, fmt.Sprintf("%s %s", "FIND_NODE", id), k)

			node.LookupContactHandleResponses(sl, id, numProbed, &channels, rpcIds)
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

func setupLookUpAlgorithm(node *Node, id *kademliaid.KademliaID) (alpha int, k int, sl *shortlist.Shortlist, channels []chan string) {
	alpha = GetEnvIntVariable("ALPHA", 3)
	k = GetEnvIntVariable("K", 5)
	sl = shortlist.NewShortlist(id, node.FindKClosest(id, nil, alpha))

	// might need more than alpha channels on final probe is closest did not change
	channels = make([]chan string, k)
	return
}

func (node *Node) LookupData(hash *kademliaid.KademliaID) string {
	alpha, k, sl, channels := setupLookUpAlgorithm(node, hash)

	// Restart the refresh timer of the bucket this ID is in range of
	bucketIndex := node.RoutingTable.GetBucketIndex(hash)
	node.RefreshTimers[bucketIndex].RestartRefreshTimer()

	// iterative lookup until the search becomes stale and no closer node
	// can be found
	result := ""
	for {
		closestSoFar := sl.Closest

		numProbed, rpcIDs := node.ProbeAlpha(sl, &channels, fmt.Sprintf("FIND_VALUE %s", hash.String()), alpha)

		if numProbed == 0 {
			log.Trace().Msg("FIND_VALUE lookup became stale")
			break
		}

		result = node.LookupDataHandleResponses(sl, hash, numProbed, &channels, rpcIDs)
		if result != "" {
			return result
		}

		if sl.Closest == closestSoFar {
			log.Trace().Msg("Closest node not updated")
			numProbed, rpcIds := node.ProbeAlpha(sl, &channels, fmt.Sprintf("%s %s", "FIND_VALUE", hash), k)

			result = node.LookupDataHandleResponses(sl, hash, numProbed, &channels, rpcIds)
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

func (node *Node) Store(value *string, contacts *[]contact.Contact, originator *contact.Contact) {
	log.Trace().Str("Value", *value).Msg("Storing value")
	node.DataStore.Insert(*value, contacts, originator, node.Network.UdpSender)
}

// FindKClosest returns a list of candidates containing the k closest nodes
// to the key being searched for (from the nodes own bucket(s))
func (node *Node) FindKClosest(key *kademliaid.KademliaID, requestorID *kademliaid.KademliaID, k int) []contact.Contact {
	return node.RoutingTable.FindClosestContacts(key, requestorID, k)
}
