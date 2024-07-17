package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	// "github.com/cosmos/cosmos-sdk/x/auth/types"
	xartypes "xarchain/x/xarchain/types"
)

// InitGenesis - Init store state from genesis data
//
// CONTRACT: old coins from the FeeCollectionKeeper need to be transferred through
// a genesis port script to the new fee collector account
func (xk Keeper) InitGenesis(ctx sdk.Context, data xartypes.GenesisState) {
	if err := xk.SetParams(ctx, data.Params); err != nil {
		panic(err)
	}

}

// ExportGenesis returns a GenesisState for a given context and keeper
// func (ak Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
// 	params := ak.GetParams(ctx)

// 	var genAccounts types.GenesisAccounts
// 	ak.IterateAccounts(ctx, func(account sdk.AccountI) bool {
// 		genAccount := account.(types.GenesisAccount)
// 		genAccounts = append(genAccounts, genAccount)
// 		return false
// 	})

// 	return types.NewGenesisState(params, genAccounts)
// }
