// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ltc

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/devwarrior777/atomicswap/libs"
	"github.com/ltcsuite/ltcd/txscript"
	"github.com/ltcsuite/ltcd/wire"
	"github.com/ltcsuite/ltcutil"
)

// auditContract pulls out information from the counterparty's contract
func auditContract(testnet bool, params libs.AuditParams) (*libs.AuditResult, error) {
	chainParams := getChainParams(testnet)

	contract, err := hex.DecodeString(params.Contract)
	if err != nil {
		return nil, fmt.Errorf("failed to decode contract: %v", err)
	}

	contractTxBytes, err := hex.DecodeString(params.ContractTx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode contract transaction: %v", err)
	}

	var contractTx wire.MsgTx
	err = contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode contract transaction: %v", err)
	}

	contractHash160 := ltcutil.Hash160(contract)
	contractOut := -1
	for i, out := range contractTx.TxOut {
		sc, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, chainParams)
		if err != nil || sc != txscript.ScriptHashTy {
			continue
		}
		if bytes.Equal(addrs[0].(*ltcutil.AddressScriptHash).Hash160()[:], contractHash160) {
			contractOut = i
			break
		}
	}
	if contractOut == -1 {
		return nil, errors.New("transaction does not contain the contract output")
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, contract)
	if err != nil {
		return nil, err
	}
	if pushes == nil {
		return nil, errors.New("contract is not an atomic swap script recognized by this tool")
	}
	if pushes.SecretSize != secretSize {
		return nil, fmt.Errorf("contract specifies strange secret size %v", pushes.SecretSize)
	}

	contractAddr, err := ltcutil.NewAddressScriptHash(contract, chainParams)
	if err != nil {
		return nil, err
	}
	recipientAddr, err := ltcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:],
		chainParams)
	if err != nil {
		return nil, err
	}
	refundAddr, err := ltcutil.NewAddressPubKeyHash(pushes.RefundHash160[:],
		chainParams)
	if err != nil {
		return nil, err
	}

	result := &libs.AuditResult{}

	result.ContractAddress = contractAddr.EncodeAddress()
	result.ContractAmount = contractTx.TxOut[contractOut].Value
	result.ContractRecipientAddress = recipientAddr.EncodeAddress()
	result.ContractRefundAddress = refundAddr.EncodeAddress()
	result.ContractRefundLocktime = pushes.LockTime
	result.ContractSecretHash = hex.EncodeToString(pushes.SecretHash[:])

	return result, nil
}
