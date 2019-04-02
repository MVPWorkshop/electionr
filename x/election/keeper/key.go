package keeper

var (
	cycleKey = []byte{0x61} // prefix for each key to an election cycle
)

// Get the key for the cycle election with primary key
func getCycleKey(primaryKey []byte) []byte {
	// Unpack primary key
	return append(cycleKey, primaryKey...)
}
