package svrtest

import (
	bnd "github.com/devwarrior777/atomicswap/libs/protobind"
)

/*
TEST DATA FOR THE LTC WALLET RPC COMMANDS

You will need your own testdata that reflects your coins configurations:
 - Testnet or not
 - RPC Info to connect to your RPC wallet node(s)
*/

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

var ltcInitiateRequest = bnd.InitiateRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
	Amount:   10000000,
}

var ltcTestnetInitiateRequest = bnd.InitiateRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
	Amount:   10000000,
}

var ltcParticipateRequest = bnd.ParticipateRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
	Amount:   10000000,
}

var ltcTestnetParticipateRequest = bnd.ParticipateRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
	Amount:   10000000,
}

var ltcRedeemRequest = bnd.RedeemRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var ltcTestnetRedeemRequest = bnd.RedeemRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var ltcRefundRequest = bnd.RefundRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var ltcTestnetRefundRequest = bnd.RefundRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var ltcPublishRequest = bnd.PublishRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var ltcTestnetPublishRequest = bnd.PublishRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var ltcGetTxRequest = bnd.GetTxRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var ltcTestnetGetTxRequest = bnd.GetTxRequest{
	Coin:     bnd.COIN_LTC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}
