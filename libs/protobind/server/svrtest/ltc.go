package svrtest

import (
	"fmt"

	bnd "github.com/devwarrior777/atomicswap/libs/protobind"
	"google.golang.org/grpc/status"
)

func testLTC(testnet bool) error {
	// LTC ping wallet
	pingreq := ltcPingWalletRPCRequest
	if testnet {
		pingreq = ltcTestnetPingWalletRPCRequest
	}
	ping, err := PingRPC(&pingreq)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf("status: %d - %v - %v", s.Code(), s.Code(), s.Message())
	}
	if ping.Errorno != bnd.ERRNO_OK {
		return fmt.Errorf("%v %s", ping.Errorno, ping.Errstr)
	}
	fmt.Println("Ping success")

	// LTC new address
	newaddressreq := ltcNewAddressRequest
	if testnet {
		newaddressreq = ltcTestnetNewAddressRequest
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
