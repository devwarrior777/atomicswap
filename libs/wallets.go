package libs

import (
	"fmt"
	"strings"

	"github.com/devwarrior777/atomicswap/libs/ltc"
	"github.com/devwarrior777/atomicswap/libs/xzc"
)

// NewAddress gets a new wallet address for the coin and network
func NewAddress(coin string, testnet bool, hostport string, rpcuser string, rpcpass string) (string, error) {
	switch strings.ToUpper(coin) {
	case "XZC":
		return newAddressXzc(testnet, hostport, rpcuser, rpcpass)
	case "LTC":
		return newAddressLtc(testnet, hostport, rpcuser, rpcpass)
	}
	//...
	return "", fmt.Errorf("unsupported coin: %s", coin)

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

//...
