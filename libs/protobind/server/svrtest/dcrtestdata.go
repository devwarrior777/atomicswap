package svrtest

import (
	bnd "github.com/devwarrior777/atomicswap/libs/protobind"
)

/*
TEST DATA FOR THE WALLET RPC COMMANDS
You will need your own testdata that reflects your coins configurations:
 - Testnet or not
 - RPC Info to connect to your DCR gRPC wallet node(s)
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

var dcrInitiateRequest = bnd.InitiateRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "",
	Amount:   10000000,
}

var dcrTestnetInitiateRequest = bnd.InitiateRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "",
	Amount:   10000000,
}

var dcrParticipateRequest = bnd.ParticipateRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "",
	Amount:   10000000,
}

var dcrTestnetParticipateRequest = bnd.ParticipateRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "",
	Amount:   10000000,
}

var dcrRedeemRequest = bnd.RedeemRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "",
}

var dcrTestnetRedeemRequest = bnd.RedeemRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "",
}

var dcrRefundRequest = bnd.RefundRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "",
}

var dcrTestnetRefundRequest = bnd.RefundRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "",
}

var dcrExtractSecretRequest = bnd.ExtractSecretRequest{
	Coin:    bnd.COIN_DCR,
	Testnet: false,
}

var dcrTestnetExtractSecretRequest = bnd.ExtractSecretRequest{
	Coin:    bnd.COIN_DCR,
	Testnet: true,
}

var dcrAuditRequest = bnd.AuditRequest{
	Coin:    bnd.COIN_DCR,
	Testnet: false,
}

var dcrTestnetAuditRequest = bnd.AuditRequest{
	Coin:    bnd.COIN_DCR,
	Testnet: true,
}

var dcrPublishRequest = bnd.PublishRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "",
}

var dcrTestnetPublishRequest = bnd.PublishRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "",
}

var dcrGetTxRequest = bnd.GetTxRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  false,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "",
}

var dcrTestnetGetTxRequest = bnd.GetTxRequest{
	Coin:     bnd.COIN_DCR,
	Testnet:  true,
	Hostport: "localhost",
	Rpcuser:  "",
	Rpcpass:  "",
	Wpass:    "123",
	Certs:    "",
}
