package arrays

import "strings"

// Convert an array of strings to an array of bytes separated by spaces.
// E.g. ["abc", "def", "geh"] -> []byte("abc def geh")
func StrArrayToByteArray(strs []string) []byte {
	return []byte(strings.Join(strs, " "))
}
