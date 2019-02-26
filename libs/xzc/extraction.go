// Copyright (c) 2017 The Decred developers
// Copyright (c) 2018 The Zcoin developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package xzc

import (
	"bytes"
	"errors"

	"github.com/zcoinofficial/xzcd/txscript"
	"github.com/zcoinofficial/xzcd/wire"
)

// extractSecret is a convenience for the participant to examine and pull out the secret from
// the initiator's redemption transaction scriptSig
func extractSecret(redemptionTx *wire.MsgTx, secretHash []byte) ([]byte, error) {
	// extractSecret loops over all pushed data from all inputs, searching for one that hashes
	// to the expected hash.  By searching through all data pushes, we avoid any
	// issues that could be caused by the initiator redeeming the participant's
	// contract with some "nonstandard" or unrecognized transaction or script
	// type.
	for _, in := range redemptionTx.TxIn {
		pushes, err := txscript.PushedData(in.SignatureScript)
		if err != nil {
			return nil, err
		}
		for _, push := range pushes {
			if bytes.Equal(sha256Hash(push), secretHash) {
				return push, nil
			}
		}
	}
	return nil, errors.New("transaction does not contain the secret")
}
