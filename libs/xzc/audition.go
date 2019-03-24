// Copyright (c) 2017/2019 The Decred developers
// Copyright (c) 2018/2019 The Zcoin developers
// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package xzc

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/zcoinofficial/xzcd/txscript"
	"github.com/zcoinofficial/xzcd/wire"
	xzcutil "github.com/zcoinofficial/xzcutil"
)

// auditContract pulls out information from the counterparty's contract
func auditContract(testnet bool, params AuditParams) (AuditResult, error) {
	result := AuditResult{}

	chainParams := getChainParams(testnet)

	contract, err := hex.DecodeString(params.Contract)
	if err != nil {
		return result, fmt.Errorf("failed to decode contract: %v", err)
	}

	contractTxBytes, err := hex.DecodeString(params.ContractTx)
	if err != nil {
		return result, fmt.Errorf("failed to decode contract transaction: %v", err)
	}

	var contractTx wire.MsgTx
	err = contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		return result, fmt.Errorf("failed to decode contract transaction: %v", err)
	}

	contractHash160 := xzcutil.Hash160(contract)
	contractOut := -1
	for i, out := range contractTx.TxOut {
		sc, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, chainParams)
		if err != nil || sc != txscript.ScriptHashTy {
			continue
		}
		if bytes.Equal(addrs[0].(*xzcutil.AddressScriptHash).Hash160()[:], contractHash160) {
			contractOut = i
			break
		}
	}
	if contractOut == -1 {
		return result, errors.New("transaction does not contain the contract output")
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(contract)
	if err != nil {
		return result, err
	}
	if pushes == nil {
		return result, errors.New("contract is not an atomic swap script recognized by this tool")
	}
	if pushes.SecretSize != secretSize {
		return result, fmt.Errorf("contract specifies strange secret size %v", pushes.SecretSize)
	}

	contractAddr, err := xzcutil.NewAddressScriptHash(contract, chainParams)
	if err != nil {
		return result, err
	}
	recipientAddr, err := xzcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:],
		chainParams)
	if err != nil {
		return result, err
	}
	refundAddr, err := xzcutil.NewAddressPubKeyHash(pushes.RefundHash160[:],
		chainParams)
	if err != nil {
		return result, err
	}

	result.ContractAddress = contractAddr.EncodeAddress()
	result.ContractAmount = contractTx.TxOut[contractOut].Value
	result.ContractRecipientAddress = recipientAddr.EncodeAddress()
	result.ContractRefundAddress = refundAddr.EncodeAddress()
	result.ContractRefundLocktime = pushes.LockTime
	result.ContractSecretHash = hex.EncodeToString(pushes.SecretHash[:])

	return result, nil
}
