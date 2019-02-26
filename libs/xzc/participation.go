// Copyright (c) 2017/2019 The Decred developers
// Copyright (c) 2018/2019 The Zcoin developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package xzc

import (
	"time"
)

// participate builds a contract & a contract transaction depending upon the hash of the
// (shared) secret. The participant will know the secret only when initiator redeems the
// contract made here
func participate(testnet bool, rpcinfo RPCInfo, params ParticipateParams) (ParticipateResult, error) {
	var result = ParticipateResult{}

	// fmt.Printf("%v\n", rpcinfo)
	// fmt.Printf("%v\n", params)

	secretHash := params.SecretHash

	// locktime after 500,000,000 (Tue Nov  5 00:53:20 1985 UTC) is interpreted
	// as a unix time rather than a block height.
	locktime := time.Now().Add(24 * time.Hour).Unix()
	// locktime := time.Now().Add(24 * time.Minute).Unix() //Test

	rpcclient, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return result, err
	}
	defer func() {
		rpcclient.Shutdown()
		rpcclient.WaitForShutdown()
	}()

	b, err := buildContract(testnet, rpcclient, &contractArgs{
		them:       params.CP1AddrP2PKH,
		amount:     params.CP1Amount,
		locktime:   locktime,
		secretHash: secretHash,
	})
	if err != nil {
		return result, err
	}

	contractFeePerKb := calcFeePerKb(b.contractFee, b.contractTx.SerializeSize())

	result.Contract = b.contract
	result.ContractP2SH = b.contractP2SH
	result.ContractTx = *b.contractTx
	result.ContractFee = b.contractFee
	result.ContractFeePerKb = contractFeePerKb

	return result, nil
}
