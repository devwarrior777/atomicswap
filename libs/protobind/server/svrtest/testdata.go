package svrtest

import (
	bnd "github.com/devwarrior777/atomicswap/libs/protobind"
)

/*
TEST DATA FOR THE WALLET RPC COMMANDS
You will need your own testdata that reflects your coins configurations:
 - Testnet or not
 - RPC Info to connect to your RPC/gRPC(DCR) wallet node(s)
 - Valid addresses for your wallet nodes
 - Valid contract, contract tx, secret, secret hash, etc. hex 'bytes',
 - Your TLS cert path for DCR gRPC wallet - if not the default
*/

/////////
// LTC //
/////////

var ltcPingWalletRPCRequest = bnd.PingWalletRPCRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var ltcTestnetPingWalletRPCRequest = bnd.PingWalletRPCRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var ltcNewAddressRequest = bnd.NewAddressRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var ltcTestnetNewAddressRequest = bnd.NewAddressRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

//...

/////////
// XZC //
/////////

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

/////////
// dcr //
/////////

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
