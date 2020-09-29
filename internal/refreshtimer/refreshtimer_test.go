package refreshtimer_test

import (
	"kademlia/internal/refreshtimer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRefreshTimer(t *testing.T) {
	bucketIndex := 20
	rt := refreshtimer.NewRefreshTimer(bucketIndex)

	// should return a refresh timer
	assert.IsType(t, refreshtimer.RefreshTimer{}, *rt)
}
