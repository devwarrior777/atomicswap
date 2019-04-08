// Copyright (c) 2017/2019 The Decred developers
// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ltc

/////////////////////////////////////////////////////////////////////
// Public command interface for the Litecoin atomic swap code library //
/////////////////////////////////////////////////////////////////////

import (
	"github.com/devwarrior777/atomicswap/libs"
)

const verify = true

const secretSize = 32

const txVersion = 2 // litecoin 0.16 needs tx v2

// PingRPC tests if wallet node RPC is available
func PingRPC(testnet bool, rpcinfo libs.RPCInfo) error {
	return pingrpc(testnet, rpcinfo)
}

// GetNewAddress gets a new address from the controlled wallet
func GetNewAddress(testnet bool, rpcinfo libs.RPCInfo) (string, error) {
	return newaddress(testnet, rpcinfo)
}

// Initiate command builds a P2SH contract and a transaction to fund it
func Initiate(testnet bool, rpcinfo libs.RPCInfo, params libs.InitiateParams) (*libs.InitiateResult, error) {
	return initiate(testnet, rpcinfo, params)
}

// Participate command builds a P2SH contract and a transaction to fund it
func Participate(testnet bool, rpcinfo libs.RPCInfo, params libs.ParticipateParams) (*libs.ParticipateResult, error) {
	return participate(testnet, rpcinfo, params)
}

// Redeem command builds a transaction to redeem a contract
func Redeem(testnet bool, rpcinfo libs.RPCInfo, params libs.RedeemParams) (*libs.RedeemResult, error) {
	return redeem(testnet, rpcinfo, params)
}

// Refund command builds a refund transaction for an unredeemed contract
func Refund(testnet bool, rpcinfo libs.RPCInfo, params libs.RefundParams) (*libs.RefundResult, error) {
	return refund(testnet, rpcinfo, params)
}

// AuditContract command
func AuditContract(testnet bool, params libs.AuditParams) (*libs.AuditResult, error) {
	return auditContract(testnet, params)
}

// Publish command broadcasts a raw hex transaction
func Publish(testnet bool, rpcinfo libs.RPCInfo, tx string) (string, error) {
	txhash, err := publish(testnet, rpcinfo, tx)
	if err != nil {
		return "", err
	}
	return txhash, nil
}

// ExtractSecret returns a secret from the scriptSig of a transaction redeeming a contract
func ExtractSecret(redemptionTx string, secretHash string) (string, error) {
	return extractSecret(redemptionTx, secretHash)
}

// GetTx gets info on a broadcasted transaction
func GetTx(testnet bool, rpcinfo libs.RPCInfo, txid string) (*libs.GetTxResult, error) {
	return getTx(testnet, rpcinfo, txid)
}

//...
