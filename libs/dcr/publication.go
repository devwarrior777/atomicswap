// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package dcr

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/decred/dcrd/chaincfg/chainhash"

	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/devwarrior777/atomicswap/libs"
)

// Publish (broadcast) transaction to the network.
func publish(testnet bool, rpcinfo libs.RPCInfo, tx string) (string, error) {
	txBytes, err := hex.DecodeString(tx)
	if err != nil {
		return "", fmt.Errorf("failed to decode broadcast transaction bytes: %v", err)
	}

	wallet, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return "", err
	}
	defer wallet.stopRPC()

	ctx := context.Background()
	response, err := wallet.client.PublishTransaction(ctx, &walletrpc.PublishTransactionRequest{
		SignedTransaction: txBytes,
	})
	if err != nil {
		return "", err
	}

	txHash, _ := chainhash.NewHash(response.TransactionHash)
	return txHash.String(), nil
}
