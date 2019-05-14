package wallets

import (
	"fmt"

	"github.com/devwarrior777/atomicswap/libs"
	bnd "github.com/devwarrior777/atomicswap/libs/protobind"
)

////////////////////////////////////
// General Wallet Node RPC Access //
////////////////////////////////////

// Wallet methods needed to access an RPC wallet node
type Wallet interface {

	// PingRPC tests if wallet node RPC is available
	PingRPC() error

	// GetNewAddress gets a new address from the controlled wallet
	GetNewAddress() (string, error)

	// Initiate command builds a P2SH contract and a transaction to fund it
	Initiate(params libs.InitiateParams) (*libs.InitiateResult, error)

	// Participate command builds a P2SH contract and a transaction to fund it
	Participate(params libs.ParticipateParams) (*libs.ParticipateResult, error)

	// Redeem command builds a transaction to redeem a contract
	Redeem(params libs.RedeemParams) (*libs.RedeemResult, error)

	// Refund command builds a refund transaction for an unredeemed contract
	Refund(params libs.RefundParams) (*libs.RefundResult, error)

	// AuditContract command
	AuditContract(params libs.AuditParams) (*libs.AuditResult, error)

	// Publish command broadcasts a raw hex transaction
	Publish(tx string) (string, error)

	// ExtractSecret returns a secret from the scriptSig of a transaction redeeming a contract
	ExtractSecret(redemptionTx string, secretHash string) (string, error)

	// GetTx gets info on a broadcasted transaction
	GetTx(txid string) (*libs.GetTxResult, error)
}

/////////////////////////////
// Individual Coin Wallets //
/////////////////////////////

// A BTCWallet can access a Bitcoin wallet node and implements Wallet
type BTCWallet struct {
	Testnet bool
	RPCInfo libs.RPCInfo
}

// An LTCWallet can access a Litecoin wallet node and implements Wallet
type LTCWallet struct {
	Testnet bool
	RPCInfo libs.RPCInfo
}

// An XZCWallet can access a Zcoin wallet node and implements Wallet
type XZCWallet struct {
	Testnet bool
	RPCInfo libs.RPCInfo
}

// A DCRWallet can access a Decred wallet and implements Wallet
type DCRWallet struct {
	Testnet bool
	RPCInfo libs.RPCInfo
}

//...

// WalletForCoin gets a concrete wallet for a coin name
// func WalletForCoin(testnet bool, rpcinfo libs.RPCInfo, coinName string) (Wallet, error) {
func WalletForCoin(testnet bool, rpcinfo libs.RPCInfo, coin bnd.COIN) (Wallet, error) {
	switch coin {
	case bnd.COIN_LTC:
		return NewLTCWallet(testnet, rpcinfo), nil
	case bnd.COIN_XZC:
		return NewXZCWallet(testnet, rpcinfo), nil
	case bnd.COIN_DCR:
		return NewDCRWallet(testnet, rpcinfo), nil
	}
	return nil, fmt.Errorf("unsupported coin %s", bnd.COIN_name[int32(coin)])
}
