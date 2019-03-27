package libs

//
// for Golang client - import this libs pkg directly
//

import (
	"fmt"

	"github.com/devwarrior777/atomicswap/libs/ltc"
	pb "github.com/devwarrior777/atomicswap/libs/protobind"
	"github.com/devwarrior777/atomicswap/libs/xzc"
)

// PingWalletRPC checks if the wallet node is available for the coin and network
func PingWalletRPC(coin pb.COIN, testnet bool, hostport string, rpcuser string, rpcpass string) error {
	switch coin {
	//case pb.COIN_BTC:
	//	return newAddressBtc(testnet, hostport, rpcuser, rpcpass)
	case pb.COIN_LTC:
		return pingWalletRPCLtc(testnet, hostport, rpcuser, rpcpass)
	case pb.COIN_XZC:
		return pingWalletRPCXzc(testnet, hostport, rpcuser, rpcpass)
	}
	return fmt.Errorf("unsupported coin: %s", coin)
}

func pingWalletRPCLtc(testnet bool, hostport string, rpcuser string, rpcpass string) error {
	rpcinfo := ltc.RPCInfo{
		HostPort: hostport,
		User:     rpcuser,
		Pass:     rpcpass,
	}
	err := ltc.PingRPC(testnet, rpcinfo)
	if err != nil {
		return err
	}
	return nil
}

func pingWalletRPCXzc(testnet bool, hostport string, rpcuser string, rpcpass string) error {
	rpcinfo := xzc.RPCInfo{
		HostPort: hostport,
		User:     rpcuser,
		Pass:     rpcpass,
	}
	err := xzc.PingRPC(testnet, rpcinfo)
	if err != nil {
		return err
	}
	return nil
}

// NewAddress gets a new wallet address for the coin and network
func NewAddress(coin pb.COIN, testnet bool, hostport string, rpcuser string, rpcpass string) (string, error) {
	switch coin {
	//case pb.COIN_BTC:
	//	return newAddressBtc(testnet, hostport, rpcuser, rpcpass)
	case pb.COIN_LTC:
		return newAddressLtc(testnet, hostport, rpcuser, rpcpass)
	case pb.COIN_XZC:
		return newAddressXzc(testnet, hostport, rpcuser, rpcpass)
	}
	return "", fmt.Errorf("unsupported coin: %s", coin)
}

func newAddressLtc(testnet bool, hostport string, rpcuser string, rpcpass string) (string, error) {
	rpcinfo := ltc.RPCInfo{
		HostPort: hostport,
		User:     rpcuser,
		Pass:     rpcpass,
	}
	newAddr, err := ltc.GetNewAddress(testnet, rpcinfo)
	if err != nil {
		return "", err
	}
	return newAddr.String(), nil
}

func newAddressXzc(testnet bool, hostport string, rpcuser string, rpcpass string) (string, error) {
	rpcinfo := xzc.RPCInfo{
		HostPort: hostport,
		User:     rpcuser,
		Pass:     rpcpass,
	}
	newAddr, err := xzc.GetNewAddress(testnet, rpcinfo)
	if err != nil {
		return "", err
	}
	return newAddr.String(), nil
}

// // Initiate creates a contract to be redeemed by the participant
// func Initiate(coin pb.COIN, testnet bool, hostport string, rpcuser string, rpcpass string, partAddress string, amount int64) (InitiateResult, error) {
// 	switch coin {
// 	//case pb.COIN_BTC:
// 	//	return newAddressBtc(testnet, hostport, rpcuser, rpcpass)
// 	case pb.COIN_LTC:
// 		return newAddressLtc(testnet, hostport, rpcuser, rpcpass, partAddress, amount)
// 	case pb.COIN_XZC:
// 		return newAddressXzc(testnet, hostport, rpcuser, rpcpass)
// 	}
// 	return InitiateResult{}, fmt.Errorf("unsupported coin: %s", coin)
// }

//...
