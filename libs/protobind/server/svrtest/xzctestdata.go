package svrtest

import (
	bnd "github.com/devwarrior777/atomicswap/libs/protobind"
)

/*
TEST DATA FOR THE XZC WALLET RPC COMMANDS
You will need your own testdata that reflects your coins configurations:
 - Testnet or not
 - RPC Info to connect to your XZC RPC wallet node(s)
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

var xzcInitiateRequest = bnd.InitiateRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
	Amount:   10000000,
}

var xzcTestnetInitiateRequest = bnd.InitiateRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
	Amount:   10000000,
}

var xzcParticipateRequest = bnd.ParticipateRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
	Amount:   10000000,
}

var xzcTestnetParticipateRequest = bnd.ParticipateRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
	Amount:   10000000,
}

var xzcRedeemRequest = bnd.RedeemRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var xzcTestnetRedeemRequest = bnd.RedeemRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var xzcRefundRequest = bnd.RefundRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var xzcTestnetRefundRequest = bnd.RefundRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var xzcExtractSecretRequest = bnd.ExtractSecretRequest{
	Coin:    bnd.COIN_XZC,
	Testnet: false,
}

var xzcTestnetExtractSecretRequest = bnd.ExtractSecretRequest{
	Coin:    bnd.COIN_XZC,
	Testnet: true,
}

var xzcAuditRequest = bnd.AuditRequest{
	Coin:    bnd.COIN_XZC,
	Testnet: false,
}

var xzcTestnetAuditRequest = bnd.AuditRequest{
	Coin:    bnd.COIN_XZC,
	Testnet: true,
}

var xzcPublishRequest = bnd.PublishRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var xzcTestnetPublishRequest = bnd.PublishRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var xzcGetTxRequest = bnd.GetTxRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}

var xzcTestnetGetTxRequest = bnd.GetTxRequest{
	Coin:     bnd.COIN_XZC,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "dev",
	Rpcpass:  "dev",
	Wpass:    "123",
	Certs:    "",
}
