package wallets

import (
	"github.com/devwarrior777/atomicswap/libs"
	"github.com/devwarrior777/atomicswap/libs/ltc"
)

// NewLTCWallet constructs an LTCWallet
func NewLTCWallet(testnet bool, rpcinfo libs.RPCInfo) *LTCWallet {
	l := &LTCWallet{
		Testnet: testnet,
		RPCInfo: rpcinfo,
	}
	return l
}

// PingRPC tests if wallet node RPC is available
func (l *LTCWallet) PingRPC() error {
	return ltc.PingRPC(l.Testnet, l.RPCInfo)
}

// GetNewAddress gets a new address from the controlled wallet
func (l *LTCWallet) GetNewAddress() (string, error) {
	return ltc.GetNewAddress(l.Testnet, l.RPCInfo)
}

// Initiate command builds a P2SH contract and a transaction to fund it
func (l *LTCWallet) Initiate(params libs.InitiateParams) (*libs.InitiateResult, error) {
	return ltc.Initiate(l.Testnet, l.RPCInfo, params)
}

// Participate command builds a P2SH contract and a transaction to fund it
func (l *LTCWallet) Participate(params libs.ParticipateParams) (*libs.ParticipateResult, error) {
	return ltc.Participate(l.Testnet, l.RPCInfo, params)
}

// Redeem command builds a transaction to redeem a contract
func (l *LTCWallet) Redeem(params libs.RedeemParams) (*libs.RedeemResult, error) {
	return ltc.Redeem(l.Testnet, l.RPCInfo, params)
}

// Refund command builds a refund transaction for an unredeemed contract
func (l *LTCWallet) Refund(params libs.RefundParams) (*libs.RefundResult, error) {
	return ltc.Refund(l.Testnet, l.RPCInfo, params)
}

// AuditContract command
func (l *LTCWallet) AuditContract(params libs.AuditParams) (*libs.AuditResult, error) {
	return ltc.AuditContract(l.Testnet, params)
}

// Publish command broadcasts a raw hex transaction
func (l *LTCWallet) Publish(tx string) (string, error) {
	return ltc.Publish(l.Testnet, l.RPCInfo, tx)
}

// ExtractSecret returns a secret from the scriptSig of a transaction redeeming a contract
func (l *LTCWallet) ExtractSecret(redemptionTx string, secretHash string) (string, error) {
	return ltc.ExtractSecret(redemptionTx, secretHash)
}

// GetTx gets info on a broadcasted transaction
func (l *LTCWallet) GetTx(txid string) (*libs.GetTxResult, error) {
	return ltc.GetTx(l.Testnet, l.RPCInfo, txid)
}
