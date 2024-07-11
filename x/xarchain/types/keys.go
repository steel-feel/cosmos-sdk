package types

const (
	// ModuleName defines the module name
	ModuleName = "xarchain"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_xarchain"
)

var (
	ParamsKey = []byte("p_xarchain")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
