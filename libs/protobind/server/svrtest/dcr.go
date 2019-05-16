package svrtest

import (
	"fmt"

	bnd "github.com/devwarrior777/atomicswap/libs/protobind"
	"google.golang.org/grpc/status"
)

func testDCR(testnet bool) error {
	// DCR ping wallet
	pingreq := dcrPingWalletRPCRequest
	if testnet {
		pingreq = dcrTestnetPingWalletRPCRequest
	}
	ping, err := pingRPC(&pingreq)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf("status: %d - %v - %v", s.Code(), s.Code(), s.Message())
	}
	if ping.Errorno != bnd.ERRNO_OK {
		return fmt.Errorf("%v %s", ping.Errorno, ping.Errstr)
	}
	fmt.Println("Ping success")

	// DCR new address
	newaddressreq := dcrNewAddressRequest
	if testnet {
		newaddressreq = dcrTestnetNewAddressRequest
	}
	newaddress, err := newAddress(&newaddressreq)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf("status: %d - %v - %v", s.Code(), s.Code(), s.Message())
	}
	if newaddress.Errorno != bnd.ERRNO_OK {
		return fmt.Errorf("%v %s", newaddress.Errorno, newaddress.Errstr)
	}
	fmt.Printf("New address: %s\n", newaddress.Address)

	return nil
}
