// Copyright (c) 2017/2019 The Decred developers
// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ltc

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ltcsuite/ltcd/txscript"
	ltcutil "github.com/ltcsuite/ltcutil"
)

// auditContract pulls out information from the counterparty's contract
func auditContract(testnet bool, params AuditParams) (AuditResult, error) {
	result := AuditResult{}

	chainParams := getChainParams(testnet)

	contract := params.Contract
	contractTx := params.ContractTx

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
		return result, errors.New("transaction does not contain the contract output")
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, contract)
	if err != nil {
		return result, err
	}
	if pushes == nil {
		return result, errors.New("contract is not an atomic swap script recognized by this tool")
	}
	if pushes.SecretSize != secretSize {
		return result, fmt.Errorf("contract specifies strange secret size %v", pushes.SecretSize)
	}

	contractAddr, err := ltcutil.NewAddressScriptHash(contract, chainParams)
	if err != nil {
		return result, err
	}
	recipientAddr, err := ltcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:],
		chainParams)
	if err != nil {
		return result, err
	}
	refundAddr, err := ltcutil.NewAddressPubKeyHash(pushes.RefundHash160[:],
		chainParams)
	if err != nil {
		return result, err
	}

	result.ContractAddress = *contractAddr
	result.ContractAmount = ltcutil.Amount(contractTx.TxOut[contractOut].Value)
	result.ContractRecipientAddress = *recipientAddr
	result.ContractRefundAddress = *refundAddr
	result.ContractRefundLocktime = pushes.LockTime
	result.ContractSecretHash = pushes.SecretHash[:]

	return result, nil
}
