package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"xarchain/x/xarchain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SyncblockAll(ctx context.Context, req *types.QueryAllSyncblockRequest) (*types.QueryAllSyncblockResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var syncblocks []types.Syncblock

    store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	syncblockStore := prefix.NewStore(store, types.KeyPrefix(types.SyncblockKeyPrefix))

	pageRes, err := query.Paginate(syncblockStore, req.Pagination, func(key []byte, value []byte) error {
		var syncblock types.Syncblock
		if err := k.cdc.Unmarshal(value, &syncblock); err != nil {
			return err
		}

		syncblocks = append(syncblocks, syncblock)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllSyncblockResponse{Syncblock: syncblocks, Pagination: pageRes}, nil
}

func (k Keeper) Syncblock(ctx context.Context, req *types.QueryGetSyncblockRequest) (*types.QueryGetSyncblockResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetSyncblock(
	    ctx,
	    req.ChainId,
        )
	if !found {
	    return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetSyncblockResponse{Syncblock: val}, nil
}