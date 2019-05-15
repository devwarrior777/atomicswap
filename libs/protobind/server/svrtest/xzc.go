package svrtest

import (
	"fmt"

	bnd "github.com/devwarrior777/atomicswap/libs/protobind"
	"google.golang.org/grpc/status"
)

func testXZC(testnet bool) error {
	// XZC ping wallet
	pingreq := xzcPingWalletRPCRequest
	if testnet {
		pingreq = xzcTestnetPingWalletRPCRequest
	}
	pingresp, err := PingRPC(&pingreq)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf("status: %d - %v - %v", s.Code(), s.Code(), s.Message())
	}
	if pingresp.Errorno != bnd.ERRNO_OK {
		return fmt.Errorf("%v %s", pingresp.Errorno, pingresp.Errstr)
	}
	fmt.Println("Ping success")

	// XZC new address
	newaddressreq := xzcNewAddressRequest
	if testnet {
		newaddressreq = xzcTestnetNewAddressRequest
	}
	newaddress, err := NewAddress(&newaddressreq)
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
