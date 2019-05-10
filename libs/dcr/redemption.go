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

	"github.com/decred/dcrd/dcrec"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/txscript"
	"github.com/decred/dcrd/wire"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/decred/dcrwallet/wallet/txrules"
	"github.com/devwarrior777/atomicswap/libs"
	// "github.com/decred/dcrd/dcrec"
	// "github.com/decred/dcrd/txscript"
	// "github.com/decred/dcrd/wire"
	// "github.com/decred/dcrutil"
	// "github.com/decred/dcrwallet/wallet/txrules"
)

// Build a transaction that can redeem the coins in the passed in contract using
// the (shared) secret
func redeem(testnet bool, rpcinfo libs.RPCInfo, params libs.RedeemParams) (*libs.RedeemResult, error) {
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

	secret, err := hex.DecodeString(params.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret: %v", err)
	}

	//--->
	pushes, err := txscript.ExtractAtomicSwapDataPushes(
		txscript.DefaultScriptVersion, contract)
	if err != nil {
		return nil, err
	}
	if pushes == nil {
		return nil, errors.New("contract is not an atomic swap script recognized by this tool")
	}
	recipientAddr, err := dcrutil.NewAddressPubKeyHash(pushes.RecipientHash160[:],
		chainParams, dcrec.STEcdsaSecp256k1)
	if err != nil {
		return nil, err
	}
	contractHash := dcrutil.Hash160(contract)
	contractOutIdx := -1
	for i, out := range contractTx.TxOut {
		sc, addrs, _, _ := txscript.ExtractPkScriptAddrs(out.Version, out.PkScript, chainParams)
		if sc == txscript.ScriptHashTy && bytes.Equal(addrs[0].Hash160()[:], contractHash) {
			contractOutIdx = i
			break
		}
	}
	if contractOutIdx == -1 {
		return nil, errors.New("transaction does not contain a contract output")
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
	fmt.Printf("redeem nar: %s\n", nar.Address)
	addr, err := dcrutil.DecodeAddress(nar.Address)
	if err != nil {
		return nil, err
	}
	outScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	contractTxHash := contractTx.TxHash()
	contractOutPoint := wire.OutPoint{
		Hash:  contractTxHash,
		Index: uint32(contractOutIdx),
		Tree:  0,
	}

	redeemTx := wire.NewMsgTx()
	redeemTx.LockTime = uint32(pushes.LockTime)
	redeemTx.AddTxIn(wire.NewTxIn(&contractOutPoint, 0, nil))
	redeemTx.AddTxOut(wire.NewTxOut(0, outScript)) // amount set below
	redeemSize := estimateRedeemSerializeSize(contract, redeemTx.TxOut)
	redeemFee := txrules.FeeForSerializeSize(feePerKb, redeemSize)
	redeemTx.TxOut[0].Value = contractTx.TxOut[contractOutIdx].Value - int64(redeemFee)
	if txrules.IsDustOutput(redeemTx.TxOut[0], feePerKb) {
		return nil, fmt.Errorf("redeem output value of %v is dust", dcrutil.Amount(redeemTx.TxOut[0].Value))
	}

	var buf bytes.Buffer
	buf.Grow(redeemTx.SerializeSize())
	redeemTx.Serialize(&buf)

	redeemSig, err := wallet.client.CreateSignature(ctx, &walletrpc.CreateSignatureRequest{
		Passphrase:            passphrase,
		Address:               recipientAddr.EncodeAddress(),
		SerializedTransaction: buf.Bytes(),
		InputIndex:            0,
		HashType:              walletrpc.CreateSignatureRequest_SIGHASH_ALL,
		PreviousPkScript:      contract,
	})
	if err != nil {
		return nil, err
	}
	redeemSigScript, err := redeemP2SHContract(contract, redeemSig.Signature,
		redeemSig.PublicKey, secret)
	if err != nil {
		return nil, err
	}
	redeemTx.TxIn[0].SignatureScript = redeemSigScript

	if verify {
		e, err := txscript.NewEngine(contractTx.TxOut[contractOutPoint.Index].PkScript,
			redeemTx, 0, verifyFlags, txscript.DefaultScriptVersion,
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

	var redeemBuf bytes.Buffer
	redeemBuf.Grow(redeemTx.SerializeSize())
	redeemTx.Serialize(&redeemBuf)
	strRefundTx := hex.EncodeToString(redeemBuf.Bytes())

	redeemFeePerKb := calcFeePerKb(redeemFee, redeemTx.SerializeSize())
	redeemTxHash := redeemTx.TxHash()
	strRedeemTxHash := redeemTxHash.String()

	var result = &libs.RedeemResult{}

	result.RedeemTx = strRefundTx
	result.RedeemTxHash = strRedeemTxHash
	result.RedeemFee = int64(redeemFee)
	result.RedeemFeePerKb = redeemFeePerKb

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
