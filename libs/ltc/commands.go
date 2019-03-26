// Copyright (c) 2017/2019 The Decred developers
// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ltc

/////////////////////////////////////////////////////////////////////
// Public command interface for the Litecoin atomic swap code library //
/////////////////////////////////////////////////////////////////////

import (
	"github.com/ltcsuite/ltcutil"
)

const verify = true

const secretSize = 32

const txVersion = 2 // litecoin 0.16 needs tx v2

// RPCInfo is RPC information passed into commands
type RPCInfo struct {
	User     string
	Pass     string
	HostPort string
}

// PingRPC tests if wallet node RPC is available
func PingRPC(testnet bool, rpcinfo RPCInfo) error {
	return pingrpc(testnet, rpcinfo)
}

// GetNewAddress gets a new address from the controlled wallet
func GetNewAddress(testnet bool, rpcinfo RPCInfo) (ltcutil.Address, error) {
	return newaddress(testnet, rpcinfo)
}

//InitiateParams is passed to the Initiate function
type InitiateParams struct {
	CP2Addr   string
	CP2Amount int64
}

//InitiateResult is returned from the Initiate function
type InitiateResult struct {
	Secret           string
	SecretHash       string
	Contract         string
	ContractP2SH     string
	ContractTx       string
	ContractTxHash   string
	ContractFee      int64
	ContractFeePerKb float64
}

// Initiate command builds a P2SH contract and a transaction to fund it
func Initiate(testnet bool, rpcinfo RPCInfo, params InitiateParams) (InitiateResult, error) {
	return initiate(testnet, rpcinfo, params)
}

//ParticipateParams is passed to the Participate command
type ParticipateParams struct {
	SecretHash string
	CP1Addr    string
	CP1Amount  int64
}

//ParticipateResult is returned from the Participate command
type ParticipateResult struct {
	Contract         string
	ContractP2SH     string
	ContractTx       string
	ContractTxHash   string
	ContractFee      int64
	ContractFeePerKb float64
}

// Participate command builds a P2SH contract and a transaction to fund it
func Participate(testnet bool, rpcinfo RPCInfo, params ParticipateParams) (ParticipateResult, error) {
	return participate(testnet, rpcinfo, params)
}

// RedeemParams is passed to the Redeem command
type RedeemParams struct {
	Secret     string
	Contract   string
	ContractTx string
}

// RedeemResult is returned from the Redeem command
type RedeemResult struct {
	RedeemTx       string
	RedeemTxHash   string
	RedeemFee      int64
	RedeemFeePerKb float64
}

// Redeem command builds a transaction to redeem a contract
func Redeem(testnet bool, rpcinfo RPCInfo, params RedeemParams) (RedeemResult, error) {
	return redeem(testnet, rpcinfo, params)
}

// RefundParams is passed to Refund command
type RefundParams struct {
	Contract   string
	ContractTx string
}

// RefundResult is returned from Refund command
type RefundResult struct {
	RefundTx       string
	RefundTxHash   string
	RefundFee      int64
	RefundFeePerKb float64
}

// Refund command builds a refund transaction for an unredeemed contract
func Refund(testnet bool, rpcinfo RPCInfo, params RefundParams) (RefundResult, error) {
	return refund(testnet, rpcinfo, params)
}

// AuditParams is passed to Audit command
type AuditParams struct {
	Contract   string
	ContractTx string
}

// AuditResult is returned from Audit command
type AuditResult struct {
	ContractAmount           int64
	ContractAddress          string
	ContractSecretHash       string
	ContractRecipientAddress string
	ContractRefundAddress    string
	ContractRefundLocktime   int64
}

// AuditContract command
func AuditContract(testnet bool, params AuditParams) (AuditResult, error) {
	return auditContract(testnet, params)
}

// Publish command broadcasts a raw hex transaction
func Publish(testnet bool, rpcinfo RPCInfo, tx string) (string, error) {
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

//...