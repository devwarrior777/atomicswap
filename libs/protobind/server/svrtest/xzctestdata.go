package svrtest

import (
	bnd "github.com/devwarrior777/atomicswap/libs/protobind"
)

/*
TEST DATA FOR THE WALLET RPC COMMANDS
You will need your own testdata that reflects your coins configurations:
 - Testnet or not
 - RPC Info to connect to your RPC/gRPC(DCR) wallet node(s)
*/

var xzcPingWalletRPCRequest = bnd.PingWalletRPCRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var xzcTestnetPingWalletRPCRequest = bnd.PingWalletRPCRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var xzcNewAddressRequest = bnd.NewAddressRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var xzcTestnetNewAddressRequest = bnd.NewAddressRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

//...
