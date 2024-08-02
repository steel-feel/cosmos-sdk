package xarchain

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"xarchain/testutil/sample"
	xarchainsimulation "xarchain/x/xarchain/simulation"
	"xarchain/x/xarchain/types"
)

// avoid unused import issue
var (
	_ = xarchainsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (


	opWeightMsgCreateIntent = "op_weight_msg_create_intent"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateIntent int = 100

	opWeightMsgUpdateIntent = "op_weight_msg_update_intent"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateIntent int = 100

	opWeightMsgCreateSyncblock = "op_weight_msg_syncblock"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateSyncblock int = 100

	opWeightMsgUpdateSyncblock = "op_weight_msg_syncblock"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateSyncblock int = 100

	opWeightMsgDeleteSyncblock = "op_weight_msg_syncblock"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteSyncblock int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	xarchainGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		SyncblockList: []types.Syncblock{
		{
			Creator: sample.AccAddress(),
ChainId: "0",
},
		{
			Creator: sample.AccAddress(),
ChainId: "1",
},
	},
	// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&xarchainGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	

	var weightMsgCreateIntent int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateIntent, &weightMsgCreateIntent, nil,
		func(_ *rand.Rand) {
			weightMsgCreateIntent = defaultWeightMsgCreateIntent
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateIntent,
		xarchainsimulation.SimulateMsgCreateIntent(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateIntent int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateIntent, &weightMsgUpdateIntent, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateIntent = defaultWeightMsgUpdateIntent
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateIntent,
		xarchainsimulation.SimulateMsgUpdateIntent(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreateSyncblock int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateSyncblock, &weightMsgCreateSyncblock, nil,
		func(_ *rand.Rand) {
			weightMsgCreateSyncblock = defaultWeightMsgCreateSyncblock
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateSyncblock,
		xarchainsimulation.SimulateMsgCreateSyncblock(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateSyncblock int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateSyncblock, &weightMsgUpdateSyncblock, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateSyncblock = defaultWeightMsgUpdateSyncblock
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateSyncblock,
		xarchainsimulation.SimulateMsgUpdateSyncblock(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteSyncblock int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteSyncblock, &weightMsgDeleteSyncblock, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteSyncblock = defaultWeightMsgDeleteSyncblock
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteSyncblock,
		xarchainsimulation.SimulateMsgDeleteSyncblock(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateIntent,
			defaultWeightMsgCreateIntent,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				xarchainsimulation.SimulateMsgCreateIntent(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateIntent,
			defaultWeightMsgUpdateIntent,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				xarchainsimulation.SimulateMsgUpdateIntent(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
	opWeightMsgCreateSyncblock,
	defaultWeightMsgCreateSyncblock,
	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		xarchainsimulation.SimulateMsgCreateSyncblock(am.accountKeeper, am.bankKeeper, am.keeper)
		return nil
	},
),
simulation.NewWeightedProposalMsg(
	opWeightMsgUpdateSyncblock,
	defaultWeightMsgUpdateSyncblock,
	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		xarchainsimulation.SimulateMsgUpdateSyncblock(am.accountKeeper, am.bankKeeper, am.keeper)
		return nil
	},
),
simulation.NewWeightedProposalMsg(
	opWeightMsgDeleteSyncblock,
	defaultWeightMsgDeleteSyncblock,
	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		xarchainsimulation.SimulateMsgDeleteSyncblock(am.accountKeeper, am.bankKeeper, am.keeper)
		return nil
	},
),
// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
