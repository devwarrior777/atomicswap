package svrtest

import (
	"fmt"
	"testing"
)

// Change these to reflect the coins you have configured
const (
	// BTC        = false
	// BTCTestnet = false
	LTC        = false
	LTCTestnet = true
	XZC        = false
	XZCTestnet = false
	DCR        = false
	DCRTestnet = true
)

func TestClient(t *testing.T) {
	// if BTC {
	// fmt.Println("\nTest BTC")
	// }
	// if BTCTestnet{
	// fmt.Println("\nTest BTC [testnet]")
	// }
	if LTC {
		fmt.Println("\nTest LTC")
		err := testLTC(false)
		if err != nil {
			t.Errorf("Status: %v\n", err.Error())
			return
		}
	}
	if LTCTestnet {
		fmt.Println("\nTest LTC [testnet]")
		err := testLTC(true)
		if err != nil {
			t.Errorf("Status: %v\n", err.Error())
			return
		}
	}
	if XZC {
		fmt.Println("\nTest XZC")
		err := testXZC(false)
		if err != nil {
			t.Errorf("Status: %v\n", err.Error())
			return
		}
	}
	if XZCTestnet {
		fmt.Println("\nTest XZC [testnet]")
		err := testXZC(true)
		if err != nil {
			t.Errorf("Status: %v\n", err.Error())
			return
		}
	}
	if DCR {
		fmt.Println("\nTest DCR")
		err := testDCR(false)
		if err != nil {
			t.Errorf("Status: %v\n", err.Error())
			return
		}
	}
	if DCRTestnet {
		fmt.Println("\nTest DCR [testnet]")
		err := testDCR(true)
		if err != nil {
			t.Errorf("Status: %v\n", err.Error())
			return
		}
	}

	//...
}
