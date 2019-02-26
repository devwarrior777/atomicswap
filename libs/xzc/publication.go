// Copyright (c) 2017/2019 The Decred developers
// Copyright (c) 2018/2019 The Zcoin developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package xzc

import (
	"github.com/zcoinofficial/xzcd/chaincfg/chainhash"
	"github.com/zcoinofficial/xzcd/wire"
)

// Publish (broadcast) transaction to the network.
func publish(testnet bool, rpcinfo RPCInfo, tx *wire.MsgTx) (*chainhash.Hash, error) {
	rpcclient, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return nil, err
	}
	defer func() {
		rpcclient.Shutdown()
		rpcclient.WaitForShutdown()
	}()

	txHash, err := sendRawTransaction(testnet, rpcclient, tx)
	if err != nil {
		return nil, err
	}

	return txHash, nil
}
