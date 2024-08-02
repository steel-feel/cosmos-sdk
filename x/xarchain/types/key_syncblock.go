package types

import "encoding/binary"

var _ binary.ByteOrder

const (
    // SyncblockKeyPrefix is the prefix to retrieve all Syncblock
	SyncblockKeyPrefix = "Syncblock/value/"
)

// SyncblockKey returns the store key to retrieve a Syncblock from the index fields
func SyncblockKey(
chainId string,
) []byte {
	var key []byte
    
    chainIdBytes := []byte(chainId)
    key = append(key, chainIdBytes...)
    key = append(key, []byte("/")...)
    
	return key
}