// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package xzc

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/devwarrior777/atomicswap/libs"
	"github.com/zcoinofficial/xzcd/chaincfg/chainhash"
	rpc "github.com/zcoinofficial/xzcd/rpcclient"
	"github.com/zcoinofficial/xzcd/txscript"
	"github.com/zcoinofficial/xzcd/wire"
	"github.com/zcoinofficial/xzcutil"
	"github.com/zcoinofficial/xzcwallet/wallet/txrules"
)

// Build a transaction that can refund the coins back to the contract creator
func refund(testnet bool, rpcinfo libs.RPCInfo, params libs.RefundParams) (libs.RefundResult, error) {
	var result = libs.RefundResult{}

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

	pushes, err := txscript.ExtractAtomicSwapDataPushes(contract)
	if err != nil {
		return result, err
	}
	if pushes == nil {
		return result, errors.New("contract is not an atomic swap script recognized by this tool")
	}

	rpcclient, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return result, err
	}
	defer func() {
		rpcclient.Shutdown()
		rpcclient.WaitForShutdown()
	}()

	feePerKb, minFeePerKb, err := getFeePerKb(rpcclient)
	if err != nil {
		return result, err
	}

	contractP2SH, err := xzcutil.NewAddressScriptHash(contract, chainParams)
	if err != nil {
		return result, err
	}
	contractP2SHPkScript, err := txscript.PayToAddrScript(contractP2SH)
	if err != nil {
		return result, err
	}

	contractTxHash := contractTx.TxHash()
	contractOutPoint := wire.OutPoint{Hash: contractTxHash, Index: ^uint32(0)}
	for i, o := range contractTx.TxOut {
		if bytes.Equal(o.PkScript, contractP2SHPkScript) {
			contractOutPoint.Index = uint32(i)
			break
		}
	}
	if contractOutPoint.Index == ^uint32(0) {
		return result, errors.New("contract tx does not contain a P2SH contract payment")
	}

	refundAddress, err := getRawChangeAddress(testnet, rpcclient)
	if err != nil {
		return result, fmt.Errorf("getrawchangeaddress: %v", err)
	}
	refundOutScript, err := txscript.PayToAddrScript(refundAddress)
	if err != nil {
		return result, err
	}

	refundAddr, err := xzcutil.NewAddressPubKeyHash(pushes.RefundHash160[:], chainParams)
	if err != nil {
		return result, err
	}

	refundTx := wire.NewMsgTx(txVersion)
	refundTx.LockTime = uint32(pushes.LockTime)
	refundTx.AddTxOut(wire.NewTxOut(0, refundOutScript)) // amount set below
	refundSize := estimateRefundSerializeSize(contract, refundTx.TxOut)
	refundFee := txrules.FeeForSerializeSize(feePerKb, refundSize)
	refundTx.TxOut[0].Value = contractTx.TxOut[contractOutPoint.Index].Value - int64(refundFee)
	if txrules.IsDustOutput(refundTx.TxOut[0], minFeePerKb) {
		return result, fmt.Errorf("refund output value of %v is dust", xzcutil.Amount(refundTx.TxOut[0].Value))
	}

	txIn := wire.NewTxIn(&contractOutPoint, nil, nil)
	txIn.Sequence = 0
	refundTx.AddTxIn(txIn)

	refundSig, refundPubKey, err := createSig(testnet, refundTx, 0, contract, refundAddr, rpcclient)
	if err != nil {
		return result, err
	}
	refundSigScript, err := refundP2SHContract(contract, refundSig, refundPubKey)
	if err != nil {
		return result, err
	}
	refundTx.TxIn[0].SignatureScript = refundSigScript

	if verify {
		e, err := txscript.NewEngine(contractTx.TxOut[contractOutPoint.Index].PkScript,
			refundTx, 0, txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(refundTx), contractTx.TxOut[contractOutPoint.Index].Value)
		if err != nil {
			return result, err
		}
		err = e.Execute()
		if err != nil {
			return result, err
		}
	}

	var refundBuf bytes.Buffer
	refundBuf.Grow(refundTx.SerializeSize())
	refundTx.Serialize(&refundBuf)
	strRefundTx := hex.EncodeToString(refundBuf.Bytes())

	var refundTxHash chainhash.Hash
	refundTxHash = refundTx.TxHash()
	strRefundTxHash := refundTxHash.String()

	result.RefundTx = strRefundTx
	result.RefundTxHash = strRefundTxHash
	result.RefundFee = int64(refundFee)
	result.RefundFeePerKb = calcFeePerKb(refundFee, refundTx.SerializeSize())

	return result, nil
}

// Build a transaction that can refund the coins back to the contract creator
func buildContractRefund(testnet bool, rpcclient *rpc.Client, contract []byte, contractTx *wire.MsgTx, feePerKb, minFeePerKb xzcutil.Amount) (refundTx *wire.MsgTx, refundFee xzcutil.Amount, err error) {
	chainParams := getChainParams(testnet)

	contractP2SH, err := xzcutil.NewAddressScriptHash(contract, chainParams)
	if err != nil {
		return nil, 0, err
	}
	contractP2SHPkScript, err := txscript.PayToAddrScript(contractP2SH)
	if err != nil {
		return nil, 0, err
	}

	contractTxHash := contractTx.TxHash()
	contractOutPoint := wire.OutPoint{Hash: contractTxHash, Index: ^uint32(0)}
	for i, o := range contractTx.TxOut {
		if bytes.Equal(o.PkScript, contractP2SHPkScript) {
			contractOutPoint.Index = uint32(i)
			break
		}
	}
	if contractOutPoint.Index == ^uint32(0) {
		return nil, 0, errors.New("contract tx does not contain a P2SH contract payment")
	}

	refundAddress, err := getRawChangeAddress(testnet, rpcclient)
	if err != nil {
		return nil, 0, fmt.Errorf("getrawchangeaddress: %v", err)
	}
	refundOutScript, err := txscript.PayToAddrScript(refundAddress)
	if err != nil {
		return nil, 0, err
	}

	pushes, _ := txscript.ExtractAtomicSwapDataPushes(contract)

	refundAddr, err := xzcutil.NewAddressPubKeyHash(pushes.RefundHash160[:], chainParams)
	if err != nil {
		return nil, 0, err
	}

	refundTx = wire.NewMsgTx(txVersion)
	refundTx.LockTime = uint32(pushes.LockTime)
	refundTx.AddTxOut(wire.NewTxOut(0, refundOutScript)) // amount set below
	refundSize := estimateRefundSerializeSize(contract, refundTx.TxOut)
	refundFee = txrules.FeeForSerializeSize(feePerKb, refundSize)
	refundTx.TxOut[0].Value = contractTx.TxOut[contractOutPoint.Index].Value - int64(refundFee)
	if txrules.IsDustOutput(refundTx.TxOut[0], minFeePerKb) {
		return nil, 0, fmt.Errorf("refund output value of %v is dust", xzcutil.Amount(refundTx.TxOut[0].Value))
	}

	txIn := wire.NewTxIn(&contractOutPoint, nil, nil)
	txIn.Sequence = 0
	refundTx.AddTxIn(txIn)

	refundSig, refundPubKey, err := createSig(testnet, refundTx, 0, contract, refundAddr, rpcclient)
	if err != nil {
		return nil, 0, err
	}
	refundSigScript, err := refundP2SHContract(contract, refundSig, refundPubKey)
	if err != nil {
		return nil, 0, err
	}
	refundTx.TxIn[0].SignatureScript = refundSigScript

	if verify {
		e, err := txscript.NewEngine(contractTx.TxOut[contractOutPoint.Index].PkScript,
			refundTx, 0, txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(refundTx), contractTx.TxOut[contractOutPoint.Index].Value)
		if err != nil {
			return nil, 0, err
		}
		err = e.Execute()
		if err != nil {
			return nil, 0, err
		}
	}

	return refundTx, refundFee, nil
}

// refundP2SHContract returns the signature script to refund a contract output
// using the contract author's signature after the locktime has been reached.
// This function assumes P2SH and appends the contract as the final data push.
func refundP2SHContract(contract, sig, pubkey []byte) ([]byte, error) {
	b := txscript.NewScriptBuilder()
	b.AddData(sig)
	b.AddData(pubkey)
	b.AddInt64(0)
	b.AddData(contract)
	return b.Script()
}
