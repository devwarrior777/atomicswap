// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ltc

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ltcsuite/ltcd/txscript"
	"github.com/ltcsuite/ltcd/wire"
)

// extractSecret is a convenience for the participant to examine and pull out the secret from
// the initiator's redemption transaction scriptSig
func extractSecret(redemptionTx string, secretHash string) (string, error) {
	// extractSecret loops over all pushed data from all inputs, searching for one that hashes
	// to the expected hash.  By searching through all data pushes, we avoid any
	// issues that could be caused by the initiator redeeming the participant's
	// contract with some "nonstandard" or unrecognized transaction or script
	// type.
	redemptionTxBytes, err := hex.DecodeString(redemptionTx)
	if err != nil {
		return "", fmt.Errorf("failed to decode redemption transaction bytes: %v", err)
	}

	var redeemTx wire.MsgTx
	err = redeemTx.Deserialize(bytes.NewReader(redemptionTxBytes))
	if err != nil {
		return "", fmt.Errorf("failed to decode redemption transaction: %v", err)
	}

	secretHashBytes, err := hex.DecodeString(secretHash)
	if err != nil {
		return "", errors.New("secret hash must be hex encoded")
	}

	if len(secretHashBytes) != sha256.Size {
		return "", errors.New("secret hash has wrong size")
	}

	for _, in := range redeemTx.TxIn {
		pushes, err := txscript.PushedData(in.SignatureScript)
		if err != nil {
			return "", err
		}
		for _, push := range pushes {
			if bytes.Equal(sha256Hash(push), secretHashBytes) {
				return hex.EncodeToString(push), nil
			}
		}
	}
	return "", errors.New("transaction does not contain the secret")
}
