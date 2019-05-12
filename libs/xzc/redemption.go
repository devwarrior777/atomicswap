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
	"github.com/zcoinofficial/xzcd/txscript"
	"github.com/zcoinofficial/xzcd/wire"
	"github.com/zcoinofficial/xzcutil"
	"github.com/zcoinofficial/xzcwallet/wallet/txrules"
)

// Build a transaction that can redeem the coins in the passed in contract using
// the (shared) secret
func redeem(testnet bool, rpcinfo libs.RPCInfo, params libs.RedeemParams) (*libs.RedeemResult, error) {
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

	secret, err := hex.DecodeString(params.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret: %v", err)
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(contract)
	if err != nil {
		return nil, err
	}
	if pushes == nil {
		return nil, errors.New("contract is not an atomic swap script recognized by this tool")
	}
	recipientAddr, err := xzcutil.NewAddressPubKeyHash(pushes.RecipientHash160[:],
		chainParams)
	if err != nil {
		return nil, err
	}
	outScript, err := txscript.PayToAddrScript(recipientAddr)
	if err != nil {
		return nil, err
	}

	contractHash := xzcutil.Hash160(contract)
	contractOutIdx := -1
	for i, out := range contractTx.TxOut {
		sc, addrs, _, _ := txscript.ExtractPkScriptAddrs(out.PkScript, chainParams)
		if sc == txscript.ScriptHashTy &&
			bytes.Equal(addrs[0].(*xzcutil.AddressScriptHash).Hash160()[:], contractHash) {
			contractOutIdx = i
			break
		}
	}
	if contractOutIdx == -1 {
		return nil, errors.New("transaction does not contain a contract output")
	}

	contractTxHash := contractTx.TxHash()
	contractOutPoint := wire.OutPoint{
		Hash:  contractTxHash,
		Index: uint32(contractOutIdx),
	}

	rpcclient, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return nil, err
	}
	defer func() {
		rpcclient.Shutdown()
		rpcclient.WaitForShutdown()
	}()

	err = walletLock(rpcclient, rpcinfo.WalletPass, 1)
	if err != nil {
		return nil, err
	}
	defer walletUnlock(rpcclient, rpcinfo.WalletPass)

	feePerKb, minFeePerKb, err := getFeePerKb(rpcclient)
	if err != nil {
		return nil, err
	}

	redeemTx := wire.NewMsgTx(txVersion)
	redeemTx.LockTime = uint32(pushes.LockTime)
	redeemTx.AddTxIn(wire.NewTxIn(&contractOutPoint, nil, nil))
	redeemTx.AddTxOut(wire.NewTxOut(0, outScript)) // amount set below
	redeemSize := estimateRedeemSerializeSize(contract, redeemTx.TxOut)
	redeemFee := txrules.FeeForSerializeSize(feePerKb, redeemSize)
	redeemTx.TxOut[0].Value = contractTx.TxOut[contractOutIdx].Value - int64(redeemFee)
	if txrules.IsDustOutput(redeemTx.TxOut[0], minFeePerKb) {
		return nil, fmt.Errorf("redeem output value of %v is dust", xzcutil.Amount(redeemTx.TxOut[0].Value))
	}

	redeemSig, redeemPubKey, err := createSig(testnet, redeemTx, 0, contract, recipientAddr, rpcclient)
	if err != nil {
		return nil, err
	}
	redeemSigScript, err := redeemP2SHContract(contract, redeemSig, redeemPubKey, secret)
	if err != nil {
		return nil, err
	}
	redeemTx.TxIn[0].SignatureScript = redeemSigScript

	if verify {
		e, err := txscript.NewEngine(contractTx.TxOut[contractOutPoint.Index].PkScript,
			redeemTx, 0, txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(redeemTx), contractTx.TxOut[contractOutIdx].Value)
		if err != nil {
			return nil, err
		}
		err = e.Execute()
		if err != nil {
			return nil, err
		}
	}

	var redeemBuf bytes.Buffer
	redeemBuf.Grow(redeemTx.SerializeSize())
	redeemTx.Serialize(&redeemBuf)
	strRefundTx := hex.EncodeToString(redeemBuf.Bytes())

	var redeemTxHash chainhash.Hash
	redeemTxHash = redeemTx.TxHash()
	strRedeemTxHash := redeemTxHash.String()

	var result = &libs.RedeemResult{}

	result.RedeemTx = strRefundTx
	result.RedeemTxHash = strRedeemTxHash
	result.RedeemFee = int64(redeemFee)
	result.RedeemFeePerKb = calcFeePerKb(redeemFee, redeemTx.SerializeSize())

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
