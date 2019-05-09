// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package dcr

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrec"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/txscript"
	"github.com/decred/dcrd/wire"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/decred/dcrwallet/wallet/txrules"
	"github.com/devwarrior777/atomicswap/libs"
)

// Build a transaction that can refund the coins back to the contract creator
func refund(testnet bool, rpcinfo libs.RPCInfo, params libs.RefundParams) (*libs.RefundResult, error) {
	chainParams := getChainParams(testnet)

	// get params suitable for dcr functions
	passphrase := []byte(rpcinfo.WalletPass)

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

	//--->
	contractP2SH, err := dcrutil.NewAddressScriptHash(contract, chainParams)
	if err != nil {
		return nil, err
	}
	contractP2SHPkScript, err := txscript.PayToAddrScript(contractP2SH)
	if err != nil {
		return nil, err
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
		return nil, errors.New("contract tx does not contain a P2SH contract payment")
	}

	wallet, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return nil, err
	}
	defer wallet.stopRPC()
	ctx := context.Background()

	nar, err := wallet.client.NextAddress(ctx, &walletrpc.NextAddressRequest{
		Account:   0, // TODO
		Kind:      walletrpc.NextAddressRequest_BIP0044_INTERNAL,
		GapPolicy: walletrpc.NextAddressRequest_GAP_POLICY_WRAP,
	})
	if err != nil {
		return nil, err
	}
	refundAddress, err := dcrutil.DecodeAddress(nar.Address)
	if err != nil {
		return nil, err
	}
	refundOutScript, err := txscript.PayToAddrScript(refundAddress)
	if err != nil {
		return nil, err
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(0, contract)
	if err != nil {
		return nil, err
	}

	refundAddr, err := dcrutil.NewAddressPubKeyHash(pushes.RefundHash160[:], chainParams,
		dcrec.STEcdsaSecp256k1)
	if err != nil {
		return nil, err
	}

	refundTx := wire.NewMsgTx()
	refundTx.LockTime = uint32(pushes.LockTime)
	refundTx.AddTxOut(wire.NewTxOut(0, refundOutScript)) // amount set below
	refundSize := estimateRefundSerializeSize(contract, refundTx.TxOut)
	refundFee := txrules.FeeForSerializeSize(feePerKb, refundSize)
	refundTx.TxOut[0].Value = contractTx.TxOut[contractOutPoint.Index].Value - int64(refundFee)
	if txrules.IsDustOutput(refundTx.TxOut[0], feePerKb) {
		return nil, fmt.Errorf("refund output value of %v is dust", dcrutil.Amount(refundTx.TxOut[0].Value))
	}

	txIn := wire.NewTxIn(&contractOutPoint, 0, nil)
	txIn.Sequence = 0
	refundTx.AddTxIn(txIn)

	var buf bytes.Buffer
	buf.Grow(refundTx.SerializeSize())
	refundTx.Serialize(&buf)

	refundSig, err := wallet.client.CreateSignature(ctx, &walletrpc.CreateSignatureRequest{
		Passphrase:            passphrase,
		Address:               refundAddr.EncodeAddress(),
		SerializedTransaction: buf.Bytes(),
		InputIndex:            0,
		HashType:              walletrpc.CreateSignatureRequest_SIGHASH_ALL,
		PreviousPkScript:      contract,
	})
	if err != nil {
		return nil, err
	}
	refundSigScript, err := refundP2SHContract(contract, refundSig.Signature,
		refundSig.PublicKey)
	if err != nil {
		return nil, err
	}
	refundTx.TxIn[0].SignatureScript = refundSigScript

	if verify {
		e, err := txscript.NewEngine(contractTx.TxOut[contractOutPoint.Index].PkScript,
			refundTx, 0, verifyFlags, txscript.DefaultScriptVersion,
			txscript.NewSigCache(10))
		if err != nil {
			return nil, err
		}
		err = e.Execute()
		if err != nil {
			return nil, err
		}
	}
	//<---

	var refundBuf bytes.Buffer
	refundBuf.Grow(refundTx.SerializeSize())
	refundTx.Serialize(&refundBuf)
	strRefundTx := hex.EncodeToString(refundBuf.Bytes())

	var refundTxHash chainhash.Hash
	refundTxHash = refundTx.TxHash()
	strRefundTxHash := refundTxHash.String()

	var result = &libs.RefundResult{}

	result.RefundTx = strRefundTx
	result.RefundTxHash = strRefundTxHash
	result.RefundFee = int64(refundFee)
	result.RefundFeePerKb = calcFeePerKb(refundFee, refundTx.SerializeSize())

	return result, nil
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
