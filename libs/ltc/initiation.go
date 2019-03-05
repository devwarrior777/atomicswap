// Copyright (c) 2017/2019 The Decred developers
// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ltc

import (
	"crypto/rand"
	"time"
)

// initiate creates a new secret then builds a contract & a contract transaction depending
// upon that secret
func initiate(testnet bool, rpcinfo RPCInfo, params InitiateParams) (InitiateResult, error) {
	var result = InitiateResult{}

	var secret32 [secretSize]byte
	_, err := rand.Read(secret32[:])
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
		them:       params.CP2AddrP2PKH,
		amount:     params.CP2Amount,
		locktime:   locktime,
		secretHash: secretHash,
	})
	if err != nil {
		return result, err
	}

	contractFeePerKb := calcFeePerKb(b.contractFee, b.contractTx.SerializeSize())

	result.Secret = secret
	result.SecretHash = secretHash
	result.Contract = b.contract
	result.ContractP2SH = b.contractP2SH
	result.ContractTx = *b.contractTx
	result.ContractFee = b.contractFee
	result.ContractFeePerKb = contractFeePerKb

	return result, nil
}
