package aura

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	aurablocklistv1 "github.com/ondoprotocol/usdy-noble/api/aura/blocklist/v1"
	aurav1 "github.com/ondoprotocol/usdy-noble/api/aura/v1"
)

func (AppModule) AutoCLIOptions() []*autocliv1.ModuleOptions {
	return []*autocliv1.ModuleOptions{
		{
			Query: &autocliv1.ServiceCommandDescriptor{
				Service: aurav1.Query_ServiceDesc.ServiceName,
				RpcCommandOptions: []*autocliv1.RpcCommandOptions{
					{
						RpcMethod: "Denom",
						Use:       "denom",
						Short:     "Query the module's denom",
					},
					{
						RpcMethod: "Paused",
						Use:       "paused",
						Short:     "Query if the module is paused",
					},
					{
						RpcMethod: "Owner",
						Use:       "owner",
						Short:     "Query the module's owner",
					},
					{
						RpcMethod: "Burners",
						Use:       "burners",
						Short:     "Query the module's burners",
					},
					{
						RpcMethod: "Minters",
						Use:       "minters",
						Short:     "Query the module's minters",
					},
					{
						RpcMethod: "Pausers",
						Use:       "pausers",
						Short:     "Query the module's pausers",
					},
					{
						RpcMethod: "BlockedChannels",
						Use:       "blocked-channels",
						Short:     "Query the blocked channels",
					},
				},
				SubCommands: map[string]*autocliv1.ServiceCommandDescriptor{
					"blocklist": {
						Service: aurablocklistv1.Query_ServiceDesc.ServiceName,
						RpcCommandOptions: []*autocliv1.RpcCommandOptions{
							{
								RpcMethod: "Owner",
								Use:       "owner",
								Short:     "Query the submodule's owner",
							},
							{
								RpcMethod: "Addresses",
								Use:       "addresses",
								Short:     "Query for all blocked addresses",
							},
							{
								RpcMethod: "Address",
								Use:       "address",
								Short:     "Query if an address is blocked",
								PositionalArgs: []*autocliv1.PositionalArgDescriptor{
									{
										ProtoField: "address",
									},
								},
							},
						},
					},
				},
			},
			Tx: &autocliv1.ServiceCommandDescriptor{
				Service: aurav1.Msg_ServiceDesc.ServiceName,
				RpcCommandOptions: []*autocliv1.RpcCommandOptions{
					{
						RpcMethod: "Burn",
						Use:       "burn [from] [amount]",
						Short:     "Transaction that burns tokens",
						PositionalArgs: []*autocliv1.PositionalArgDescriptor{
							{
								ProtoField: "from",
							},
							{
								ProtoField: "amount",
							},
						},
					},
					{
						RpcMethod: "Mint",
						Use:       "mint [to] [amount]",
						Short:     "Transaction that mints tokens",
						PositionalArgs: []*autocliv1.PositionalArgDescriptor{
							{
								ProtoField: "to",
							},
							{
								ProtoField: "amount",
							},
						},
					},
					{
						RpcMethod: "Pause",
						Use:       "pause",
						Short:     "Transaction that pauses the module",
					},
					{
						RpcMethod: "Unpause",
						Use:       "unpause",
						Short:     "Transaction that unpauses the module",
					},
					{
						RpcMethod: "TransferOwnership",
						Use:       "transfer-ownership [new-owner]",
						Short:     "Transfer ownership of the module",
						PositionalArgs: []*autocliv1.PositionalArgDescriptor{
							{
								ProtoField: "new_owner",
							},
						},
					},
					{
						RpcMethod: "AcceptOwnership",
						Use:       "accept-ownership",
						Short:     "Accept ownership of the module",
						Long:      "Accept ownership of the module, assuming there is an pending ownership transfer",
					},
					{
						RpcMethod: "AddBurner",
						Use:       "add-burner [burner] [allowance]",
						Short:     "Add a new burner with an initial allowance",
						PositionalArgs: []*autocliv1.PositionalArgDescriptor{
							{
								ProtoField: "burner",
							},
							{
								ProtoField: "allowance",
							},
						},
					},
					{
						RpcMethod: "RemoveBurner",
						Use:       "remove-burner [burner]",
						Short:     "Remove an existing burner",
						PositionalArgs: []*autocliv1.PositionalArgDescriptor{
							{
								ProtoField: "burner",
							},
						},
					},
					{
						RpcMethod: "SetBurnerAllowance",
						Use:       "set-burner-allowance [burner] [allowance]",
						Short:     "Set an existing burner's allowance",
						PositionalArgs: []*autocliv1.PositionalArgDescriptor{
							{
								ProtoField: "burner",
							},
							{
								ProtoField: "allowance",
							},
						},
					},
					{
						RpcMethod: "AddMinter",
						Use:       "add-minter [minter] [allowance]",
						Short:     "Add a new minter with an initial allowance",
						PositionalArgs: []*autocliv1.PositionalArgDescriptor{
							{
								ProtoField: "minter",
							},
							{
								ProtoField: "allowance",
							},
						},
					},
					{
						RpcMethod: "RemoveMinter",
						Use:       "remove-minter [minter]",
						Short:     "Remove an existing minter",
						PositionalArgs: []*autocliv1.PositionalArgDescriptor{
							{
								ProtoField: "minter",
							},
						},
					},
					{
						RpcMethod: "SetMinterAllowance",
						Use:       "set-minter-allowance [minter] [allowance]",
						Short:     "Set an existing minter's allowance",
						PositionalArgs: []*autocliv1.PositionalArgDescriptor{
							{
								ProtoField: "minter",
							},
							{
								ProtoField: "allowance",
							},
						},
					},
					{
						RpcMethod: "AddPauser",
						Use:       "add-pauser [pauser]",
						Short:     "Add a new pauser",
						PositionalArgs: []*autocliv1.PositionalArgDescriptor{
							{
								ProtoField: "pauser",
							},
						},
					},
					{
						RpcMethod: "RemovePauser",
						Use:       "remove-pauser [pauser]",
						Short:     "Remove an existing pauser",
						PositionalArgs: []*autocliv1.PositionalArgDescriptor{
							{
								ProtoField: "pauser",
							},
						},
					},
					{
						RpcMethod: "AddBlockedChannel",
						Use:       "add-blocked-channel [channel]",
						Short:     "Add a new blocked channel",
						PositionalArgs: []*autocliv1.PositionalArgDescriptor{
							{
								ProtoField: "channel",
							},
						},
					},
					{
						RpcMethod: "RemoveBlockedChannel",
						Use:       "remove-blocked-channel [channel]",
						Short:     "Remove an existing blocked channel",
						PositionalArgs: []*autocliv1.PositionalArgDescriptor{
							{
								ProtoField: "channel",
							},
						},
					},
				},
				SubCommands: map[string]*autocliv1.ServiceCommandDescriptor{
					"blocklist": {
						Service: aurablocklistv1.Msg_ServiceDesc.ServiceName,
						RpcCommandOptions: []*autocliv1.RpcCommandOptions{
							{
								RpcMethod: "TransferOwnership",
								Use:       "transfer-ownership [new-owner]",
								Short:     "Transfer ownership of submodule",
								PositionalArgs: []*autocliv1.PositionalArgDescriptor{
									{
										ProtoField: "new_owner",
									},
								},
							},
							{
								RpcMethod: "AcceptOwnership",
								Use:       "accept-ownership",
								Short:     "Accept ownership of submodule",
								Long:      "Accept ownership of submodule, assuming there is an pending ownership transfer",
							},
							{
								RpcMethod: "AddToBlocklist",
								Use:       "add-to-blocklist [addresses ...]",
								Short:     "Add addresses to the blocklist",
								PositionalArgs: []*autocliv1.PositionalArgDescriptor{
									{
										ProtoField: "accounts",
									},
								},
							},
							{
								RpcMethod: "RemoveFromBlocklist",
								Use:       "remove-from-blocklist [addresses ...]",
								Short:     "Remove addresses from the blocklist",
								PositionalArgs: []*autocliv1.PositionalArgDescriptor{
									{
										ProtoField: "accounts",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
