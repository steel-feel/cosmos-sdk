package keeper

import (
	"context"

	"xarchain/x/xarchain/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetCblock set cblock in the store
func (k Keeper) SetCblock(ctx context.Context, cblock types.Cblock) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.CblockKey))
	b := k.cdc.MustMarshal(&cblock)
	store.Set([]byte{0}, b)
}

// GetCblock returns cblock
func (k Keeper) GetCblock(ctx context.Context) (val types.Cblock, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.CblockKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveCblock removes cblock from the store
func (k Keeper) RemoveCblock(ctx context.Context) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.CblockKey))
	store.Delete([]byte{0})
}
