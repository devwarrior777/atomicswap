// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package dcr

import (
	"bytes"
	"context"
	"fmt"

	"github.com/decred/dcrwallet/rpc/walletrpc"

	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/txscript"
	"github.com/decred/dcrd/wire"
	"golang.org/x/crypto/ripemd160"
)

// contractArgs specifies the common parameters used to create the initiator's
// and participant's contract.
type contractArgs struct {
	them       *dcrutil.AddressPubKeyHash
	amount     dcrutil.Amount
	locktime   int64
	secretHash []byte
}

// builtContract houses the details regarding a contract and the contract
// payment transaction, as well as the transaction to perform a refund.
type builtContract struct {
	contract     []byte
	contractP2SH dcrutil.Address
	contractTx   wire.MsgTx
	contractFee  dcrutil.Amount
}

// buildContract creates a contract for the parameters specified in args, using
// wallet RPC to generate an internal address to redeem the refund and to sign
// the payment to the contract transaction.
func buildContract(testnet bool, c walletrpc.WalletServiceClient, args *contractArgs, p string) (*builtContract, error) {
	ctx := context.Background()
	passphrase := []byte(p)
	chainParams := getChainParams(testnet)

	nar, err := c.NextAddress(ctx, &walletrpc.NextAddressRequest{
		Account:   0, // TODO
		Kind:      walletrpc.NextAddressRequest_BIP0044_INTERNAL,
		GapPolicy: walletrpc.NextAddressRequest_GAP_POLICY_WRAP,
	})
	if err != nil {
		return nil, err
	}
	refundAddr, err := dcrutil.DecodeAddress(nar.Address)
	if err != nil {
		return nil, err
	}
	if _, ok := refundAddr.(*dcrutil.AddressPubKeyHash); !ok {
		return nil, fmt.Errorf("NextAddress: address %v is not P2PKH", refundAddr)
	}

	contract, err := atomicSwapContract(refundAddr.Hash160(), args.them.Hash160(),
		args.locktime, args.secretHash)
	if err != nil {
		return nil, err
	}
	contractP2SH, err := dcrutil.NewAddressScriptHash(contract, chainParams)
	if err != nil {
		return nil, err
	}
	contractP2SHPkScript, err := txscript.PayToAddrScript(contractP2SH)
	if err != nil {
		return nil, err
	}

	ctr, err := c.ConstructTransaction(ctx, &walletrpc.ConstructTransactionRequest{
		SourceAccount: 0, // TODO
		NonChangeOutputs: []*walletrpc.ConstructTransactionRequest_Output{{
			Destination: &walletrpc.ConstructTransactionRequest_OutputDestination{
				Script:        contractP2SHPkScript,
				ScriptVersion: 0,
			},
			Amount: int64(args.amount),
		}},
	})
	if err != nil {
		return nil, err
	}
	contractFee := dcrutil.Amount(ctr.TotalPreviousOutputAmount - ctr.TotalOutputAmount)
	str, err := c.SignTransaction(ctx, &walletrpc.SignTransactionRequest{
		Passphrase:            passphrase,
		SerializedTransaction: ctr.UnsignedTransaction,
	})
	if err != nil {
		return nil, err
	}
	var contractTx wire.MsgTx
	err = contractTx.Deserialize(bytes.NewReader(str.Transaction))
	if err != nil {
		return nil, err
	}

	return &builtContract{
		contract,
		contractP2SH,
		contractTx,
		contractFee,
	}, nil
}

// atomicSwapContract returns an output script that may be redeemed by one of
// two signature scripts:
//
//   <their sig> <their pubkey> <initiator secret> 1
//
//   <my sig> <my pubkey> 0
//
// The first signature script is the normal redemption path done by the other
// party and requires the initiator's secret.  The second signature script is
// the refund path performed by us, but the refund can only be performed after
// locktime.
func atomicSwapContract(pkhMe, pkhThem *[ripemd160.Size]byte, locktime int64, secretHash []byte) ([]byte, error) {
	b := txscript.NewScriptBuilder()

	b.AddOp(txscript.OP_IF) // Normal redeem path
	{
		// Require initiator's secret to be a known length that the redeeming
		// party can audit.  This is used to prevent fraud attacks between two
		// currencies that have different maximum data sizes.
		b.AddOp(txscript.OP_SIZE)
		b.AddInt64(secretSize)
		b.AddOp(txscript.OP_EQUALVERIFY)

		// Require initiator's secret to be known to redeem the output.
		b.AddOp(txscript.OP_SHA256)
		b.AddData(secretHash)
		b.AddOp(txscript.OP_EQUALVERIFY)

		// Verify their signature is being used to redeem the output.  This
		// would normally end with OP_EQUALVERIFY OP_CHECKSIG but this has been
		// moved outside of the branch to save a couple bytes.
		b.AddOp(txscript.OP_DUP)
		b.AddOp(txscript.OP_HASH160)
		b.AddData(pkhThem[:])
	}
	b.AddOp(txscript.OP_ELSE) // Refund path
	{
		// Verify locktime and drop it off the stack (which is not done by
		// CLTV).
		b.AddInt64(locktime)
		b.AddOp(txscript.OP_CHECKLOCKTIMEVERIFY)
		b.AddOp(txscript.OP_DROP)

		// Verify our signature is being used to redeem the output.  This would
		// normally end with OP_EQUALVERIFY OP_CHECKSIG but this has been moved
		// outside of the branch to save a couple bytes.
		b.AddOp(txscript.OP_DUP)
		b.AddOp(txscript.OP_HASH160)
		b.AddData(pkhMe[:])
	}
	b.AddOp(txscript.OP_ENDIF)

	// Complete the signature check.
	b.AddOp(txscript.OP_EQUALVERIFY)
	b.AddOp(txscript.OP_CHECKSIG)

	return b.Script()
}
