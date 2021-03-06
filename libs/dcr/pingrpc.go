package dcr

// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

import (
	"github.com/devwarrior777/atomicswap/libs"
)

// pingrpc tests if wallet node RPC is available
func pingrpc(testnet bool, rpcinfo libs.RPCInfo) error {
	wallet, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return err
	}
	defer wallet.stopRPC()

	err = wallet.ping()
	if err != nil {
		return err
	}

	return nil
}
