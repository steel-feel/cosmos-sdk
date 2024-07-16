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
					RpcMethod:      "ShowTask",
					Use:            "show-task [id]",
					Short:          "Query show-task",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},

				{
					RpcMethod:      "ListTask",
					Use:            "list-task",
					Short:          "Query list-task",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},

				{
					RpcMethod: "Cblock",
					Use:       "show-cblock",
					Short:     "show cblock",
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
					RpcMethod:      "CreateTask",
					Use:            "create-task [title]",
					Short:          "Send a create-task tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "title"}},
				},
				{
					RpcMethod:      "UpdateTask",
					Use:            "update-task [title] [status] [id]",
					Short:          "Send a update-task tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "title"}, {ProtoField: "status"}, {ProtoField: "id"}},
				},
				{
					RpcMethod:      "CreateCblock",
					Use:            "create-cblock [blocknumber]",
					Short:          "Create cblock",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "blocknumber"}},
				},
				{
					RpcMethod:      "UpdateCblock",
					Use:            "update-cblock [blocknumber]",
					Short:          "Update cblock",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "blocknumber"}},
				},
				{
					RpcMethod: "DeleteCblock",
					Use:       "delete-cblock",
					Short:     "Delete cblock",
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
