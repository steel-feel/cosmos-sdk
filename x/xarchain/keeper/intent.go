package keeper

import (
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"xarchain/x/xarchain/types"
)

func (k Keeper) SetIntent(ctx sdk.Context, intent types.Intent) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.IntentKey))
	b := k.cdc.MustMarshal(&intent)
	store.Set([]byte(intent.Intentid), b)
}

func (k Keeper) AppendIntent(ctx sdk.Context, intent types.Intent) uint64 {
	count := k.GetIntentCount(ctx)
	intent.Id = count
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.IntentKey))
	appendedValue := k.cdc.MustMarshal(&intent)
	store.Set([]byte(intent.Intentid), appendedValue)
	k.SetIntentCount(ctx, count+1)

	//Emit the intent event, 
	if err := ctx.EventManager().EmitTypedEvent(&types.IntentEvent{
		Intentid: intent.Intentid,
		Chainid: intent.ChainId,
		Txhash: intent.Txhash,
	}); err != nil {
		return 0
	}

	return count
}

func (k Keeper) SetIntentCount(ctx sdk.Context, count uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.IntentCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

func (k Keeper) GetIntentById(ctx sdk.Context, intentid string) (val types.Intent, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.IntentKey))
	b := store.Get([]byte(intentid))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
/*
func (k Keeper) GetIntentByIntentHash(ctx sdk.Context, intentHash string) (val types.Intent, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.IntentKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.IntentKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var intent types.Intent
		k.cdc.MustUnmarshal(iterator.Value(), &intent)
		if intent.IntentHash == intentHash {
			return intent, true
		}
	}
	return val, false
}*/

func GetPostIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

func GetIntentIdBytes(hash string) []byte {
	return []byte(hash)
}

func (k Keeper) GetIntentCount(ctx sdk.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.IntentCountKey)
	bz := store.Get(byteKey)
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

