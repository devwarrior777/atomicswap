package wallets

import (
	"github.com/devwarrior777/atomicswap/libs"
	"github.com/devwarrior777/atomicswap/libs/dcr"
)

// NewDCRWallet constructs an DCRWallet
func NewDCRWallet(testnet bool, rpcinfo libs.RPCInfo) *DCRWallet {
	d := &DCRWallet{
		Testnet: testnet,
		RPCInfo: rpcinfo,
	}
	return d
}

// PingRPC tests if wallet node RPC is available
func (d *DCRWallet) PingRPC() error {
	return dcr.PingRPC(d.Testnet, d.RPCInfo)
}

// GetNewAddress gets a new address from the controlled wallet
func (d *DCRWallet) GetNewAddress() (string, error) {
	return dcr.GetNewAddress(d.Testnet, d.RPCInfo)
}

// Initiate command builds a P2SH contract and a transaction to fund it
func (d *DCRWallet) Initiate(params libs.InitiateParams) (*libs.InitiateResult, error) {
	return dcr.Initiate(d.Testnet, d.RPCInfo, params)
}

// Participate command builds a P2SH contract and a transaction to fund it
func (d *DCRWallet) Participate(params libs.ParticipateParams) (*libs.ParticipateResult, error) {
	return dcr.Participate(d.Testnet, d.RPCInfo, params)
}

// Redeem command builds a transaction to redeem a contract
func (d *DCRWallet) Redeem(params libs.RedeemParams) (*libs.RedeemResult, error) {
	return dcr.Redeem(d.Testnet, d.RPCInfo, params)
}

// Refund command builds a refund transaction for an unredeemed contract
func (d *DCRWallet) Refund(params libs.RefundParams) (*libs.RefundResult, error) {
	return dcr.Refund(d.Testnet, d.RPCInfo, params)
}

// AuditContract command
func (d *DCRWallet) AuditContract(params libs.AuditParams) (*libs.AuditResult, error) {
	return dcr.AuditContract(d.Testnet, params)
}

// Publish command broadcasts a raw hex transaction
func (d *DCRWallet) Publish(tx string) (string, error) {
	return dcr.Publish(d.Testnet, d.RPCInfo, tx)
}

// ExtractSecret returns a secret from the scriptSig of a transaction redeeming a contract
func (d *DCRWallet) ExtractSecret(redemptionTx string, secretHash string) (string, error) {
	return dcr.ExtractSecret(redemptionTx, secretHash)
}

// GetTx gets info on a broadcasted transaction
func (d *DCRWallet) GetTx(txid string) (*libs.GetTxResult, error) {
	return dcr.GetTx(d.Testnet, d.RPCInfo, txid)
}
