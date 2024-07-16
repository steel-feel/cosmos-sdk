package xarchain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"xarchain/x/xarchain/keeper"
	"xarchain/x/xarchain/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set if defined
	if genState.Cblock != nil {
		k.SetCblock(ctx, *genState.Cblock)
	}
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Get all cblock
	cblock, found := k.GetCblock(ctx)
	if found {
		genesis.Cblock = &cblock
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
