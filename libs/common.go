package libs

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/rand"
	"time"
)

const hexstr32len = 32 * 2

// These are the common structures used with all the swap coins
// for making transactions and other wallet functionality for
// atomic swaps

// RPCInfo is RPC information passed into commands
// HostPort:	If no  port specified defaults to the coin's default
// 				port for the network
type RPCInfo struct {
	User       string // RPC Username
	Pass       string // RPC Password
	HostPort   string // RPC host[:port] can be ipv4 [ipv6]
	WalletPass string // Wallet-passphrase
	Certs      string // DCR Wallet
}

//InitiateParams is passed to the Initiate function
type InitiateParams struct {
	SecretHash string // Hash of the shared secret
	CP2Addr    string // Counterparty 2 (Participant) Adddress
	CP2Amount  int64  // Amount (sats) to pay into Participant redeemable contract
}

//InitiateResult is returned from the Initiate function
type InitiateResult struct {
	Contract         string
	ContractP2SH     string
	ContractTx       string
	ContractTxHash   string
	ContractFee      int64
	ContractFeePerKb float64
}

//ParticipateParams is passed to the Participate command
type ParticipateParams struct {
	SecretHash string
	CP1Addr    string // Counterparty 1 (Initiator) contract Adddress
	CP1Amount  int64  // Amount (sats) to pay into Initiator redeemable contract
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

// GetTxResult is returned from GetTx command
type GetTxResult struct {
	Confirmations uint64
	Blockhash     string
	Blockindex    int
	Blocktime     uint64
	Time          uint64
	TimeReceived  uint64
	Hex           string
}

// GetRand32 creates a 32-'byte' pseudo random hex string
func GetRand32() string {
	src := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 32)
	_, _ = src.Read(b)
	return hex.EncodeToString(b)[:]
}

// Hash256 takes a 32-'byte' hex string and hashes the binary number
// represented then outputs as a hex string
func Hash256(s string) (string, error) {
	if len(s) != hexstr32len {
		return "", errors.New("hash256 - invalid input")
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}
	if len(b) != 32 {
		return "", errors.New("hash256 - invalid input byte length")
	}
	h32 := sha256.Sum256(b)
	return hex.EncodeToString(h32[:]), nil
}
