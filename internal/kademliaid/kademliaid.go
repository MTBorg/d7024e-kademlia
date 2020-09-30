package kademliaid

import (
	"crypto/sha1"
	"encoding/hex"
	"math"
	"math/rand"
	"time"
)

// the static number of bytes in a KademliaID
const IDLength = 20

// type definition of a KademliaID
type KademliaID [IDLength]byte

// NewKademliaID returns a new instance of a KademliaID based on the string input
func NewKademliaID(data *string) KademliaID {
	hash := sha1.Sum([]byte(*data))
	return hash
}

// NewRandomKademliaID returns a new instance of a random KademliaID,
// change this to a better version if you like
func NewRandomKademliaID() *KademliaID {
	rand.Seed(time.Now().UTC().UnixNano()) // Update the seed
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

// Less returns true if kademliaID < otherKademliaID (bitwise)
func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

// Equals returns true if kademliaID == otherKademliaID (bitwise)
func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

// CalcDistance returns a new instance of a KademliaID that is built
// through a bitwise XOR operation betweeen kademliaID and target
func (kademliaID KademliaID) CalcDistance(target *KademliaID) *KademliaID {
	result := KademliaID{}
	for i := 0; i < IDLength; i++ {
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}

// String returns a simple string representation of a KademliaID
func (kademliaID *KademliaID) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}

func FromString(s string) *KademliaID {
	id := KademliaID{}
	decoded, _ := hex.DecodeString(s)
	copy(id[:], decoded)
	return &id
}

// NewKademliaIDInRange returns a new random KademliaID which will be
// inside the range of the nodes (with NodeID id) k-bucket with index bucketIndex
//
// The implementation assumes that the first byte in the byte array of the
// kademliaID is the most significant byte (big endian)
func NewKademliaIDInRange(id *KademliaID, bucketIndex int) *KademliaID {
	// The routing tables GetBucketIndex will say that an id which differs on all
	// bits is in bucket with index 0 so the order of the buckets is reversed
	// from the description in the paper, hence the convertion below
	bucketIndex = 159 - bucketIndex
	commonID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		// Copy the node id until the differing bit is reached, this is bit nr
		// bucketIndex + 1 (assuming bit starts at 0), randomize the remaining
		var iByte uint8 = 0
		for j := 0; j < 8; j++ {
			bit := uint8(math.Pow(float64(2), float64(7-j)))
			bitValue := id[i] & bit

			if (IDLength-i)*8-j > bucketIndex+1 {
				// Copy the bits if still at the matching prefix
				iByte += bitValue
			} else if (IDLength-i)*8-j == bucketIndex+1 {
				// make sure the differing bit differs
				if bitValue == 0 {
					iByte += bit
				}
			} else {
				// Randomize all the bits after the differing bit
				iByte += uint8(rand.Intn(2)) * uint8(math.Pow(float64(2), float64(7-j)))
			}
			commonID[i] = iByte
		}
	}
	return &commonID
}
