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
	"github.com/ltcsuite/ltcd/wire"
	ltcutil "github.com/ltcsuite/ltcutil"
	"github.com/ltcsuite/ltcwallet/wallet/txrules"
)

// Build a transaction that can redeem the coins in the passed in contract using
// the (shared) secret
func redeem(testnet bool, rpcinfo RPCInfo, params RedeemParams) (RedeemResult, error) {
	var result = RedeemResult{}

	chainParams := getChainParams(testnet)

	secret := params.Secret
	contract := params.Contract
	contractTx := params.ContractTx

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, contract)
	if err != nil {
		return result, err
	}
	if pushes == nil {
		return result, errors.New("contract is not an atomic swap script recognized by this tool")
	}
	recipientAddr, err := ltcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:],
		chainParams)
	if err != nil {
		return result, err
	}
	contractHash := ltcutil.Hash160(contract)
	contractOut := -1
	for i, out := range contractTx.TxOut {
		sc, addrs, _, _ := txscript.ExtractPkScriptAddrs(out.PkScript, chainParams)
		if sc == txscript.ScriptHashTy &&
			bytes.Equal(addrs[0].(*ltcutil.AddressScriptHash).Hash160()[:], contractHash) {
			contractOut = i
			break
		}
	}
	if contractOut == -1 {
		return result, errors.New("transaction does not contain a contract output")
	}

	rpcclient, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return result, err
	}
	defer func() {
		rpcclient.Shutdown()
		rpcclient.WaitForShutdown()
	}()

	addr, err := getRawChangeAddress(testnet, rpcclient)
	if err != nil {
		return result, fmt.Errorf("getrawchangeaddres: %v", err)
	}
	outScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return result, err
	}

	contractTxHash := contractTx.TxHash()
	contractOutPoint := wire.OutPoint{
		Hash:  contractTxHash,
		Index: uint32(contractOut),
	}

	feePerKb, minFeePerKb, err := getFeePerKb(rpcclient)
	if err != nil {
		return result, err
	}

	redeemTx := wire.NewMsgTx(txVersion)
	redeemTx.LockTime = uint32(pushes.LockTime)
	redeemTx.AddTxIn(wire.NewTxIn(&contractOutPoint, nil, nil))
	redeemTx.AddTxOut(wire.NewTxOut(0, outScript)) // amount set below
	redeemSize := estimateRedeemSerializeSize(contract, redeemTx.TxOut)
	fee := txrules.FeeForSerializeSize(feePerKb, redeemSize)
	redeemTx.TxOut[0].Value = contractTx.TxOut[contractOut].Value - int64(fee)
	if txrules.IsDustOutput(redeemTx.TxOut[0], minFeePerKb) {
		return result, fmt.Errorf("redeem output value of %v is dust", ltcutil.Amount(redeemTx.TxOut[0].Value))
	}

	redeemSig, redeemPubKey, err := createSig(testnet, redeemTx, 0, contract, recipientAddr, rpcclient)
	if err != nil {
		return result, err
	}
	redeemSigScript, err := redeemP2SHContract(contract, redeemSig, redeemPubKey, secret)
	if err != nil {
		return result, err
	}
	redeemTx.TxIn[0].SignatureScript = redeemSigScript

	if verify {
		e, err := txscript.NewEngine(contractTx.TxOut[contractOutPoint.Index].PkScript,
			redeemTx, 0, txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(redeemTx), contractTx.TxOut[contractOut].Value)
		if err != nil {
			return result, err
		}
		err = e.Execute()
		if err != nil {
			return result, err
		}
	}

	result.RedeemTx = *redeemTx
	result.RedeemFee = fee
	result.RedeemFeePerKb = calcFeePerKb(fee, redeemTx.SerializeSize())

	return result, nil
}

// redeemP2SHContract returns the signature script to redeem a contract output
// using the redeemer's signature and the initiator's secret.  This function
// assumes P2SH and appends the contract as the final data push.
func redeemP2SHContract(contract, sig, pubkey, secret []byte) ([]byte, error) {
	b := txscript.NewScriptBuilder()
	b.AddData(sig)
	b.AddData(pubkey)
	b.AddData(secret)
	b.AddInt64(1)
	b.AddData(contract)
	return b.Script()
}
