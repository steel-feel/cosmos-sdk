package types

const (
	// ModuleName defines the module name
	ModuleName = "xarchain"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_xarchain"

	// taskKey is used to uniquely identify tasks within the system.
	// It will be used as the beginning of the key for each task, followed bei their unique ID
	TaskKey = "Task/value/"

	// This key will be used to keep track of the ID of the latest Task added to the store.
	TaskCountKey = "Task/count/"

		// taskKey is used to uniquely identify tasks within the system.
	// It will be used as the beginning of the key for each task, followed bei their unique ID
	IntentKey = "Intent/value/"

	// This key will be used to keep track of the ID of the latest Task added to the store.
	IntentCountKey = "Intent/count/"
)

var (
	ParamsKey = []byte("p_xarchain")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	CblockKey = "Cblock/value/"
)
