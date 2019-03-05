// Copyright (c) 2017/2019 The Decred developers
// Copyright (c) 2018/2019 The Zcoin developers
// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package xzc

import (
	xzcutil "github.com/zcoinofficial/xzcutil"
)

// getNewAddress gets a new wallet address from the controlled wallet
func newaddress(testnet bool, rpcinfo RPCInfo) (xzcutil.Address, error) {
	rpcclient, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return nil, err
	}
	defer func() {
		rpcclient.Shutdown()
		rpcclient.WaitForShutdown()
	}()

	addr, err := getNewAddress(testnet, rpcclient)
	if err != nil {
		return nil, err
	}

	return addr, nil
}
