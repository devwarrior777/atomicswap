package main

// THIS WILL BE REPLACED WITH BETTER I/F

//
// For Golang client - ignore this file as we do not need a server.
// import the libs & libs/<COIN> packages directly
//
// For node, python, etc. make client bindings for your language
// using the atomicswap.proto as the definitions
//

import (
	"fmt"

	"github.com/devwarrior777/atomicswap/libs"
	bnd "github.com/devwarrior777/atomicswap/libs/protobind"

	"github.com/devwarrior777/atomicswap/libs/ltc"
	"github.com/devwarrior777/atomicswap/libs/xzc"
	//...
)

// PingWalletRPC checks if the wallet node is available for the coin and network
// The server wants an error returned but we embed our errors in the binding. If
// the response gets to the client with an error it means there was a transport error
func pingWalletRPC(request *bnd.PingWalletRPCRequest) *bnd.PingWalletRPCResponse {
	switch request.Coin {
	//case pb.COIN_BTC:
	//	return newAddressBtc(testnet, hostport, rpcuser, rpcpass)
	case bnd.COIN_LTC:
		return pingWalletRPCLtc(request)
	case bnd.COIN_XZC:
		return pingWalletRPCXzc(request)
		//
		//...more coins
	}
	response := &bnd.PingWalletRPCResponse{Errorno: bnd.ERRNO_UNSUPPORTED}
	response.Errstr = fmt.Sprintf("unsupported coin: %v", request.Coin)
	return response
}

func pingWalletRPCLtc(request *bnd.PingWalletRPCRequest) *bnd.PingWalletRPCResponse {
	response := &bnd.PingWalletRPCResponse{}
	rpcinfo := libs.RPCInfo{
		HostPort: request.Hostport,
		User:     request.Rpcuser,
		Pass:     request.Rpcpass,
	}
	err := ltc.PingRPC(request.Testnet, rpcinfo)
	if err != nil {
		response.Errorno = bnd.ERRNO_LIBS
		response.Errstr = err.Error()
		return response
	}

	response.Errorno = bnd.ERRNO_OK
	return response
}

func pingWalletRPCXzc(request *bnd.PingWalletRPCRequest) *bnd.PingWalletRPCResponse {
	response := &bnd.PingWalletRPCResponse{}
	rpcinfo := libs.RPCInfo{
		HostPort: request.Hostport,
		User:     request.Rpcuser,
		Pass:     request.Rpcpass,
	}
	err := xzc.PingRPC(request.Testnet, rpcinfo)
	if err != nil {
		response.Errorno = bnd.ERRNO_LIBS
		response.Errstr = err.Error()
		return response
	}

	response.Errorno = bnd.ERRNO_OK
	return response
}

//...
