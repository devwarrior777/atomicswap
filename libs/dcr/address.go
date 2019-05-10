// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package dcr

import (
	"context"

	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/devwarrior777/atomicswap/libs"
)

// newaddress gets a new wallet address from the controlled wallet
func newaddress(testnet bool, rpcinfo libs.RPCInfo) (string, error) {
	wallet, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return "", err
	}
	defer wallet.stopRPC()
	ctx := context.Background()

	nar, err := wallet.client.NextAddress(ctx, &walletrpc.NextAddressRequest{
		Account:   0, // TODO
		Kind:      walletrpc.NextAddressRequest_BIP0044_INTERNAL,
		GapPolicy: walletrpc.NextAddressRequest_GAP_POLICY_WRAP,
	})
	if err != nil {
		return "", err
	}

	return nar.Address, nil
}
