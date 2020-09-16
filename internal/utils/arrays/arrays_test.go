package arrays_test

import (
	"kademlia/internal/utils/arrays"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrArrayToByteArray(t *testing.T) {
	//should return correct result
	assert.Equal(t, arrays.StrArrayToByteArray([]string{"abc", "def", "geh"}), []byte("abc def geh"))
}
