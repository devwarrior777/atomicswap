// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package dcr

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/devwarrior777/atomicswap/libs"
	// "github.com/jrick/btcd/chaincfg/chainhash"
	"github.com/decred/dcrd/chaincfg/chainhash"

	"github.com/decred/dcrd/dcrutil"
)

// initiate builds a contract & a contract transaction depending on the secret hash parameter
// passed in
func initiate(testnet bool, rpcinfo libs.RPCInfo, params libs.InitiateParams) (*libs.InitiateResult, error) {
	chainParams := getChainParams(testnet)

	cp2Addr, err := dcrutil.DecodeAddress(params.CP2Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode participant address: %v", err)
	}
	if !cp2Addr.IsForNet(chainParams) {
		return nil, fmt.Errorf("participant address is not "+
			"intended for use on %v", chainParams.Name)
	}

	cp2AddrP2PKH, ok := cp2Addr.(*dcrutil.AddressPubKeyHash)
	if !ok {
		return nil, errors.New("participant address is not P2PKH")
	}

	cp2Amount := dcrutil.Amount(params.CP2Amount)

	secretHash, err := hex.DecodeString(params.SecretHash)
	if err != nil {
		return nil, errors.New("secret hash must be hex encoded")
	}
	if len(secretHash) != sha256.Size {
		return nil, errors.New("secret hash has wrong size")
	}

	// locktime after 500,000,000 (Tue Nov  5 00:53:20 1985 UTC) is interpreted
	// as a unix time rather than a block height.
	locktime := time.Now().Add(48 * time.Hour).Unix()
	// locktime := time.Now().Add(48 * time.Minute).Unix() //Test

	wallet, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return nil, err
	}
	defer wallet.stopRPC()

	b, err := buildContract(testnet, wallet.client, &contractArgs{
		them:       cp2AddrP2PKH,
		amount:     cp2Amount,
		locktime:   locktime,
		secretHash: secretHash,
	}, rpcinfo.WalletPass)
	if err != nil {
		return nil, err
	}

	//->
	//fmt.Printf("%v %v %v %v \n\n%v\n", cp2AddrP2PKH, cp2Addr, cp2Amount, locktime, b)
	//return nil, errors.New("debug")
	//<-

	contractFeePerKb := calcFeePerKb(b.contractFee, b.contractTx.SerializeSize())

	var contractBuf bytes.Buffer
	contractBuf.Grow(b.contractTx.SerializeSize())
	b.contractTx.Serialize(&contractBuf)
	strContractTx := hex.EncodeToString(contractBuf.Bytes())

	var contractTxHash chainhash.Hash
	contractTxHash = b.contractTx.TxHash()
	strContractTxHash := contractTxHash.String()

	var result = &libs.InitiateResult{}

	result.Contract = hex.EncodeToString(b.contract)
	result.ContractP2SH = b.contractP2SH.EncodeAddress()
	result.ContractTx = strContractTx
	result.ContractTxHash = strContractTxHash
	result.ContractFee = int64(b.contractFee)
	result.ContractFeePerKb = contractFeePerKb
	result.ContractRefundLocktime = locktime

	return result, nil
}
