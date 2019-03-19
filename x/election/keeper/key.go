package keeper

import "github.com/MVPWorkshop/legaler-bc/x/election/types"

var (
	cycleKey = []byte{0x61} // prefix for each key to an election cycle
)

// Get the key for the cycle election with primary key
func getCycleKey(primaryKey types.Hash) []byte {
	// Convert primary key to byte slice, and unpack it
	return append(cycleKey, primaryKey[:]...)
}
