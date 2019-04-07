// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ltc

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/devwarrior777/atomicswap/libs"
	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	"github.com/ltcsuite/ltcutil"
)

// initiate creates a new secret then builds a contract & a contract transaction depending
// upon that secret
func initiate(testnet bool, rpcinfo libs.RPCInfo, params libs.InitiateParams) (libs.InitiateResult, error) {
	var result = libs.InitiateResult{}

	chainParams := getChainParams(testnet)

	cp2Addr, err := ltcutil.DecodeAddress(params.CP2Addr, chainParams)
	if err != nil {
		return result, fmt.Errorf("failed to decode participant address: %v", err)
	}
	if !cp2Addr.IsForNet(chainParams) {
		return result, fmt.Errorf("participant address is not "+
			"intended for use on %v", chainParams.Name)
	}

	cp2AddrP2PKH, ok := cp2Addr.(*ltcutil.AddressPubKeyHash)
	if !ok {
		return result, errors.New("participant address is not P2PKH")
	}

	cp2Amount := ltcutil.Amount(params.CP2Amount)

	var secret32 [secretSize]byte
	_, err = rand.Read(secret32[:])
	if err != nil {
		return result, err
	}
	secret := make([]byte, secretSize)
	copy(secret, secret32[:])
	secretHash := sha256Hash(secret[:])

	// locktime after 500,000,000 (Tue Nov  5 00:53:20 1985 UTC) is interpreted
	// as a unix time rather than a block height.
	locktime := time.Now().Add(48 * time.Hour).Unix()
	// locktime := time.Now().Add(48 * time.Minute).Unix() //Test

	rpcclient, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return result, err
	}
	defer func() {
		rpcclient.Shutdown()
		rpcclient.WaitForShutdown()
	}()

	b, err := buildContract(testnet, rpcclient, &contractArgs{
		them:       cp2AddrP2PKH,
		amount:     cp2Amount,
		locktime:   locktime,
		secretHash: secretHash,
	})
	if err != nil {
		return result, err
	}

	contractFeePerKb := calcFeePerKb(b.contractFee, b.contractTx.SerializeSize())

	var contractBuf bytes.Buffer
	contractBuf.Grow(b.contractTx.SerializeSize())
	b.contractTx.Serialize(&contractBuf)
	strContractTx := hex.EncodeToString(contractBuf.Bytes())

	var contractTxHash chainhash.Hash
	contractTxHash = b.contractTx.TxHash()
	strContractTxHash := contractTxHash.String()

	result.Secret = hex.EncodeToString(secret)
	result.SecretHash = hex.EncodeToString(secretHash)
	result.Contract = hex.EncodeToString(b.contract)
	result.ContractP2SH = b.contractP2SH.EncodeAddress()
	result.ContractTx = strContractTx
	result.ContractTxHash = strContractTxHash
	result.ContractFee = int64(b.contractFee)
	result.ContractFeePerKb = contractFeePerKb

	return result, nil
}
