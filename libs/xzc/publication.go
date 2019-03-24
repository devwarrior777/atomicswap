// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package xzc

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/zcoinofficial/xzcd/wire"
)

// Publish (broadcast) transaction to the network.
func publish(testnet bool, rpcinfo RPCInfo, tx string) (string, error) {
	txBytes, err := hex.DecodeString(tx)
	if err != nil {
		return "", fmt.Errorf("failed to decode broadcast transaction bytes: %v", err)
	}

	var broadcastTx wire.MsgTx
	err = broadcastTx.Deserialize(bytes.NewReader(txBytes))
	if err != nil {
		return "", fmt.Errorf("failed to decode broadcast transaction: %v", err)
	}

	rpcclient, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return "", err
	}
	defer func() {
		rpcclient.Shutdown()
		rpcclient.WaitForShutdown()
	}()

	txHash, err := sendRawTransaction(testnet, rpcclient, &broadcastTx)
	if err != nil {
		return "", err
	}

	return txHash.String(), nil
}
