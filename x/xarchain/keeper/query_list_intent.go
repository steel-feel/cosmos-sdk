package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
"fmt"
	"xarchain/x/xarchain/types"
)

func (k Keeper) ListIntent(ctx context.Context, req *types.QueryListIntentRequest) (*types.QueryListIntentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.IntentKey))

	var intents []types.Intent
	
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var intent types.Intent
		if err := k.cdc.Unmarshal(value, &intent); err != nil {
			return err
		}

		intents = append(intents, intent)
		return nil
	})


	fmt.Printf("Intents: %v\n", intents)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}




	return &types.QueryListIntentResponse{Intent: intents, Pagination: pageRes}, nil
}
