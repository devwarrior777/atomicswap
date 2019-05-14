package wallets

import (
	"github.com/devwarrior777/atomicswap/libs"
	"github.com/devwarrior777/atomicswap/libs/xzc"
)

// NewXZCWallet constructs an XZCWallet
func NewXZCWallet(testnet bool, rpcinfo libs.RPCInfo) *XZCWallet {
	x := &XZCWallet{
		Testnet: testnet,
		RPCInfo: rpcinfo,
	}
	return x
}

// PingRPC tests if wallet node RPC is available
func (x *XZCWallet) PingRPC() error {
	return xzc.PingRPC(x.Testnet, x.RPCInfo)
}

// GetNewAddress gets a new address from the controlled wallet
func (x *XZCWallet) GetNewAddress() (string, error) {
	return xzc.GetNewAddress(x.Testnet, x.RPCInfo)
}

// Initiate command builds a P2SH contract and a transaction to fund it
func (x *XZCWallet) Initiate(params libs.InitiateParams) (*libs.InitiateResult, error) {
	return xzc.Initiate(x.Testnet, x.RPCInfo, params)
}

// Participate command builds a P2SH contract and a transaction to fund it
func (x *XZCWallet) Participate(params libs.ParticipateParams) (*libs.ParticipateResult, error) {
	return xzc.Participate(x.Testnet, x.RPCInfo, params)
}

// Redeem command builds a transaction to redeem a contract
func (x *XZCWallet) Redeem(params libs.RedeemParams) (*libs.RedeemResult, error) {
	return xzc.Redeem(x.Testnet, x.RPCInfo, params)
}

// Refund command builds a refund transaction for an unredeemed contract
func (x *XZCWallet) Refund(params libs.RefundParams) (*libs.RefundResult, error) {
	return xzc.Refund(x.Testnet, x.RPCInfo, params)
}

// AuditContract command
func (x *XZCWallet) AuditContract(params libs.AuditParams) (*libs.AuditResult, error) {
	return xzc.AuditContract(x.Testnet, params)
}

// Publish command broadcasts a raw hex transaction
func (x *XZCWallet) Publish(tx string) (string, error) {
	return xzc.Publish(x.Testnet, x.RPCInfo, tx)
}

// ExtractSecret returns a secret from the scriptSig of a transaction redeeming a contract
func (x *XZCWallet) ExtractSecret(redemptionTx string, secretHash string) (string, error) {
	return xzc.ExtractSecret(redemptionTx, secretHash)
}

// GetTx gets info on a broadcasted transaction
func (x *XZCWallet) GetTx(txid string) (*libs.GetTxResult, error) {
	return xzc.GetTx(x.Testnet, x.RPCInfo, txid)
}
