package xarchain

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "xarchain/api/xarchain/xarchain"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},


				{
					RpcMethod:      "GetIntent",
					Use:            "get-intent [id]",
					Short:          "Query get-intent",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},

				{
					RpcMethod:      "ListIntent",
					Use:            "list-intent",
					Short:          "Query list-intent",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},

				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
			
				{
					RpcMethod:      "CreateIntent",
					Use:            "create-intent [from] [to] [data] [value] [chain-id]",
					Short:          "Send a create-intent tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "from"}, {ProtoField: "to"}, {ProtoField: "data"}, {ProtoField: "value"}, {ProtoField: "chainId"}},
				},
				{
					RpcMethod:      "UpdateIntent",
					Use:            "update-intent [status] [id]",
					Short:          "Send a update-intent tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "status"}, {ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
