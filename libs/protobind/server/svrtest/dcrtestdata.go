package svrtest

import (
	bnd "github.com/devwarrior777/atomicswap/libs/protobind"
)

/*
TEST DATA FOR THE WALLET RPC COMMANDS
You will need your own testdata that reflects your coins configurations:
 - Testnet or not
 - RPC Info to connect to your RPC/gRPC(DCR) wallet node(s)
 - Your TLS cert path for DCR gRPC wallet - if not the default
*/

var dcrPingWalletRPCRequest = bnd.PingWalletRPCRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "", // default cert path
}

var dcrTestnetPingWalletRPCRequest = bnd.PingWalletRPCRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "", // default cert path
}

var dcrNewAddressRequest = bnd.NewAddressRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "", // default cert path
}

var dcrTestnetNewAddressRequest = bnd.NewAddressRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "", // default cert path
}

//...
