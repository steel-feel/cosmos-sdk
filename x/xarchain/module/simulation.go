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
	opWeightMsgCreateCblock = "op_weight_msg_cblock"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateCblock int = 100

	opWeightMsgUpdateCblock = "op_weight_msg_cblock"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateCblock int = 100

	opWeightMsgDeleteCblock = "op_weight_msg_cblock"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteCblock int = 100

	opWeightMsgCreateIntent = "op_weight_msg_create_intent"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateIntent int = 100

	opWeightMsgUpdateIntent = "op_weight_msg_update_intent"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateIntent int = 100

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
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&xarchainGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateCblock int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateCblock, &weightMsgCreateCblock, nil,
		func(_ *rand.Rand) {
			weightMsgCreateCblock = defaultWeightMsgCreateCblock
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateCblock,
		xarchainsimulation.SimulateMsgCreateCblock(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateCblock int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateCblock, &weightMsgUpdateCblock, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateCblock = defaultWeightMsgUpdateCblock
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateCblock,
		xarchainsimulation.SimulateMsgUpdateCblock(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteCblock int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteCblock, &weightMsgDeleteCblock, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteCblock = defaultWeightMsgDeleteCblock
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteCblock,
		xarchainsimulation.SimulateMsgDeleteCblock(am.accountKeeper, am.bankKeeper, am.keeper),
	))

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

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateCblock,
			defaultWeightMsgCreateCblock,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				xarchainsimulation.SimulateMsgCreateCblock(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateCblock,
			defaultWeightMsgUpdateCblock,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				xarchainsimulation.SimulateMsgUpdateCblock(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteCblock,
			defaultWeightMsgDeleteCblock,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				xarchainsimulation.SimulateMsgDeleteCblock(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
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
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
