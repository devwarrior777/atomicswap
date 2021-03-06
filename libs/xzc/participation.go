// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package xzc

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/devwarrior777/atomicswap/libs"

	"github.com/zcoinofficial/xzcd/chaincfg/chainhash"
	"github.com/zcoinofficial/xzcutil"
)

// participate builds a contract & a contract transaction depending upon the hash of the
// (shared) secret. The participant will know the secret only when initiator redeems the
// contract made here
func participate(testnet bool, rpcinfo libs.RPCInfo, params libs.ParticipateParams) (*libs.ParticipateResult, error) {
	chainParams := getChainParams(testnet)

	cp1Addr, err := xzcutil.DecodeAddress(params.CP1Addr, chainParams)
	if err != nil {
		return nil, fmt.Errorf("failed to decode initiator address: %v", err)
	}
	if !cp1Addr.IsForNet(chainParams) {
		return nil, fmt.Errorf("initiator address is not intended for use on %v", chainParams.Name)
	}

	cp1Address, ok := cp1Addr.(*xzcutil.AddressPubKeyHash)
	if !ok {
		return nil, errors.New("initiator address is not P2PKH")
	}

	cp1Amount := xzcutil.Amount(params.CP1Amount)

	secretHashBytes, err := hex.DecodeString(params.SecretHash)
	if err != nil {
		return nil, errors.New("secret hash must be hex encoded")
	}
	if len(secretHashBytes) != sha256.Size {
		return nil, errors.New("secret hash has wrong size")
	}

	// locktime after 500,000,000 (Tue Nov  5 00:53:20 1985 UTC) is interpreted
	// as a unix time rather than a block height.
	locktime := time.Now().Add(24 * time.Hour).Unix()
	// locktime := time.Now().Add(24 * time.Minute).Unix() //Test

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

	b, err := buildContract(testnet, rpcclient, &contractArgs{
		them:       cp1Address,
		amount:     cp1Amount,
		locktime:   locktime,
		secretHash: secretHashBytes,
	})
	if err != nil {
		return nil, err
	}

	contractFeePerKb := calcFeePerKb(b.contractFee, b.contractTx.SerializeSize())

	var contractBuf bytes.Buffer
	contractBuf.Grow(b.contractTx.SerializeSize())
	b.contractTx.Serialize(&contractBuf)
	strContractTx := hex.EncodeToString(contractBuf.Bytes())

	var contractTxHash chainhash.Hash
	contractTxHash = b.contractTx.TxHash()
	strContractTxHash := contractTxHash.String()

	var result = &libs.ParticipateResult{}

	result.Contract = hex.EncodeToString(b.contract)
	result.ContractP2SH = b.contractP2SH.EncodeAddress()
	result.ContractTx = strContractTx
	result.ContractTxHash = strContractTxHash
	result.ContractFee = int64(b.contractFee)
	result.ContractFeePerKb = contractFeePerKb
	result.ContractRefundLocktime = locktime

	return result, nil
}
