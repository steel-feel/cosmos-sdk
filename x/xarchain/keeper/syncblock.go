package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"xarchain/x/xarchain/types"
)

// SetSyncblock set a specific syncblock in the store from its index
func (k Keeper) SetSyncblock(ctx context.Context, syncblock types.Syncblock) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store :=  prefix.NewStore(storeAdapter, types.KeyPrefix(types.SyncblockKeyPrefix))
	b := k.cdc.MustMarshal(&syncblock)
	store.Set(types.SyncblockKey(
        syncblock.ChainId,
    ), b)
}

// GetSyncblock returns a syncblock from its index
func (k Keeper) GetSyncblock(
    ctx context.Context,
    chainId string,
    
) (val types.Syncblock, found bool) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SyncblockKeyPrefix))

	b := store.Get(types.SyncblockKey(
        chainId,
    ))
    if b == nil {
        return val, false
    }

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveSyncblock removes a syncblock from the store
func (k Keeper) RemoveSyncblock(
    ctx context.Context,
    chainId string,
    
) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SyncblockKeyPrefix))
	store.Delete(types.SyncblockKey(
	    chainId,
    ))
}

// GetAllSyncblock returns all syncblock
func (k Keeper) GetAllSyncblock(ctx context.Context) (list []types.Syncblock) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
    store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.SyncblockKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Syncblock
		k.cdc.MustUnmarshal(iterator.Value(), &val)
        list = append(list, val)
	}

    return
}
