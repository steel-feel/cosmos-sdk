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
	opWeightMsgCreateTask = "op_weight_msg_create_task"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateTask int = 100

	opWeightMsgUpdateTask = "op_weight_msg_update_task"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateTask int = 100

	opWeightMsgCreateCblock = "op_weight_msg_cblock"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateCblock int = 100

	opWeightMsgUpdateCblock = "op_weight_msg_cblock"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateCblock int = 100

	opWeightMsgDeleteCblock = "op_weight_msg_cblock"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteCblock int = 100

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

	var weightMsgCreateTask int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateTask, &weightMsgCreateTask, nil,
		func(_ *rand.Rand) {
			weightMsgCreateTask = defaultWeightMsgCreateTask
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateTask,
		xarchainsimulation.SimulateMsgCreateTask(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateTask int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateTask, &weightMsgUpdateTask, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateTask = defaultWeightMsgUpdateTask
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateTask,
		xarchainsimulation.SimulateMsgUpdateTask(am.accountKeeper, am.bankKeeper, am.keeper),
	))

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

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreateTask,
			defaultWeightMsgCreateTask,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				xarchainsimulation.SimulateMsgCreateTask(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateTask,
			defaultWeightMsgUpdateTask,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				xarchainsimulation.SimulateMsgUpdateTask(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
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
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
