// Copyright (c) 2017/2019 The Decred developers
// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/devwarrior777/atomicswap/libs"
	"github.com/devwarrior777/atomicswap/libs/xzc" // Use new libs/xzc pkg
	"github.com/zcoinofficial/xzcd/txscript"
)

var (
	flagset     = flag.NewFlagSet("", flag.ExitOnError)
	connectFlag = flagset.String("s", "localhost", "host[:port] of Zcoin Core wallet RPC server")
	rpcuserFlag = flagset.String("rpcuser", "", "username for wallet RPC authentication")
	rpcpassFlag = flagset.String("rpcpass", "", "password for wallet RPC authentication")
	testnetFlag = flagset.Bool("testnet", false, "use testnet network")
	walletPass  = flagset.String("wpass", "", "wallet passphrase")
)

// There are two directions that the atomic swap can be performed, as the
// initiator can be on either chain.  This tool only deals with creating the
// Zcoin transactions for these swaps.  A second tool should be used for the
// transaction on the other chain.  Any chain can be used so long as it supports
// OP_SHA256 and OP_CHECKLOCKTIMEVERIFY.
//
// Example scenerios using zcoin as the second chain:
//
// Scenerio 1:
//   cp1 initiates (dcr)
//   cp2 participates with cp1 H(S) (xzc)
//   cp1 redeems xzc revealing S
//     - must verify H(S) in contract is hash of known secret
//   cp2 redeems dcr with S
//
// Scenerio 2:
//   cp1 initiates (xzc)
//   cp2 participates with cp1 H(S) (dcr)
//   cp1 redeems dcr revealing S
//     - must verify H(S) in contract is hash of known secret
//   cp2 redeems xzc with S

func init() {
	flagset.Usage = func() {
		fmt.Println("Usage: xzcatomicswap [flags] cmd [cmd args]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  initiate <participant address> <amount>")
		fmt.Println("  participate <initiator address> <amount> <secret hash>")
		fmt.Println("  redeem <contract> <contract transaction> <secret>")
		fmt.Println("  refund <contract> <contract transaction>")
		fmt.Println("  extractsecret <redemption transaction> <secret hash>")
		fmt.Println("  auditcontract <contract> <contract transaction>")
		fmt.Println("  gettx <txid>")
		fmt.Println("  newaddress")
		fmt.Println()
		fmt.Println("Flags:")
		flagset.PrintDefaults()
	}
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if err != nil {
		os.Exit(1)
	}
}

func checkCmdArgLength(args []string, required int) (nArgs int) {
	if len(args) < required {
		return 0
	}
	for i, arg := range args[:required] {
		if len(arg) != 1 && strings.HasPrefix(arg, "-") {
			return i
		}
	}
	return required
}

func run() error {
	flagset.Parse(os.Args[1:])
	args := flagset.Args()
	if len(args) == 0 {
		flagset.Usage()
		return errors.New("no args")
	}
	cmdArgs := 0
	switch args[0] {
	case "initiate":
		cmdArgs = 2
	case "participate":
		cmdArgs = 3
	case "redeem":
		cmdArgs = 3
	case "refund":
		cmdArgs = 2
	case "extractsecret":
		cmdArgs = 2
	case "auditcontract":
		cmdArgs = 2
	case "gettx":
		cmdArgs = 1
	case "newaddress":
		cmdArgs = 0
	default:
		flagset.Usage()
		return fmt.Errorf("unknown command %v", args[0])
	}
	nArgs := checkCmdArgLength(args[1:], cmdArgs)
	flagset.Parse(args[1+nArgs:])
	if nArgs < cmdArgs {
		flagset.Usage()
		return fmt.Errorf("%s: too few arguments", args[0])
	}
	if flagset.NArg() != 0 {
		flagset.Usage()
		return fmt.Errorf("unexpected argument: %s", flagset.Arg(0))
	}

	switch args[0] {
	case "initiate":
		return initiate(args)

	case "participate":
		return participate(args)

	case "redeem":
		return redeem(args)

	case "refund":
		return refund(args)

	case "extractsecret":
		return extractSecret(args)

	case "auditcontract":
		return auditContract(args)

	case "gettx":
		return getTx(args)

	case "newaddress":
		return newAddress(args)
	}
	flagset.Usage()
	return fmt.Errorf("unexpected argument: %s", flagset.Arg(0))
}

func initiate(args []string) error {
	amountF64, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return fmt.Errorf("failed to decode amount: %v", err)
	}

	amount, err := xzc.NewAmount(amountF64)
	if err != nil {
		return err
	}

	var rpcinfo libs.RPCInfo
	rpcinfo.HostPort = *connectFlag
	rpcinfo.User = *rpcuserFlag
	rpcinfo.Pass = *rpcpassFlag
	rpcinfo.WalletPass = *walletPass

	err = xzc.PingRPC(*testnetFlag, rpcinfo)
	if err != nil {
		return fmt.Errorf("Ping RPC: error: %v", err)
	}

	secret := libs.GetRand32()
	secretHash, err := libs.Hash256(secret)
	if err != nil {
		return errors.New("cannot generate a secret hash")
	}

	var params libs.InitiateParams
	params.SecretHash = secretHash
	params.CP2Addr = args[1]
	params.CP2Amount = int64(amount)

	var result *libs.InitiateResult
	result, err = xzc.Initiate(*testnetFlag, rpcinfo, params)
	if err != nil {
		return fmt.Errorf("Initiate: %v", err)
	}

	fmt.Printf("Secret:      %s\n", secret)
	fmt.Printf("Secret hash: %s\n\n", secretHash)
	fmt.Printf("Contract fee: %d (%0.8f XZC/kB)\n", result.ContractFee, result.ContractFeePerKb)
	fmt.Printf("Contract (%s):\n", result.ContractP2SH)
	fmt.Printf("%s\n\n", result.Contract)
	fmt.Printf("Contract transaction (%s):\n", result.ContractTxHash)
	fmt.Printf("%s\n\n", result.ContractTx)

	doPublish, err := askPublishTx("contract")
	if err != nil {
		return err
	}
	if doPublish {
		txHash, err := xzc.Publish(*testnetFlag, rpcinfo, result.ContractTx)
		if err != nil {
			return err
		}
		fmt.Printf("Published %s transaction (%s)\n", "contract", txHash)
	}

	return nil
}

func participate(args []string) error {
	amountF64, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return fmt.Errorf("failed to decode amount: %v", err)
	}
	amount, err := xzc.NewAmount(amountF64)
	if err != nil {
		return err
	}

	var rpcinfo libs.RPCInfo
	rpcinfo.HostPort = *connectFlag
	rpcinfo.User = *rpcuserFlag
	rpcinfo.Pass = *rpcpassFlag
	rpcinfo.WalletPass = *walletPass

	err = xzc.PingRPC(*testnetFlag, rpcinfo)
	if err != nil {
		return fmt.Errorf("Ping RPC: error: %v", err)
	}

	var params libs.ParticipateParams
	params.SecretHash = args[3]
	params.CP1Addr = args[1]
	params.CP1Amount = int64(amount)

	var result *libs.ParticipateResult
	result, err = xzc.Participate(*testnetFlag, rpcinfo, params)
	if err != nil {
		return fmt.Errorf("Participate: %v", err)
	}

	fmt.Printf("Contract fee: %d (%0.8f XZC/kB)\n", result.ContractFee, result.ContractFeePerKb)
	fmt.Printf("Contract (%s):\n", result.ContractP2SH)
	fmt.Printf("%s\n\n", result.Contract)
	fmt.Printf("Contract transaction (%s):\n", result.ContractTxHash)
	fmt.Printf("%s\n\n", result.ContractTx)

	doPublish, err := askPublishTx("contract")
	if err != nil {
		return err
	}
	if doPublish {
		txHash, err := xzc.Publish(*testnetFlag, rpcinfo, result.ContractTx)
		if err != nil {
			return err
		}
		fmt.Printf("Published %s transaction (%s)\n", "contract", txHash)
	}

	return nil
}

func redeem(args []string) error {
	var rpcinfo libs.RPCInfo
	rpcinfo.HostPort = *connectFlag
	rpcinfo.User = *rpcuserFlag
	rpcinfo.Pass = *rpcpassFlag
	rpcinfo.WalletPass = *walletPass

	err := xzc.PingRPC(*testnetFlag, rpcinfo)
	if err != nil {
		return fmt.Errorf("Ping RPC: error: %v", err)
	}

	var params libs.RedeemParams
	params.Contract = args[1]
	params.ContractTx = args[2]
	params.Secret = args[3]

	var result *libs.RedeemResult
	result, err = xzc.Redeem(*testnetFlag, rpcinfo, params)
	if err != nil {
		return fmt.Errorf("Redeem: %v", err)
	}

	fmt.Printf("Redeem fee:   %d (%0.8f XZC/kB)\n\n", result.RedeemFee, result.RedeemFeePerKb)
	fmt.Printf("Redeem transaction (%s):\n", result.RedeemTxHash)
	fmt.Printf("%s\n\n", result.RedeemTx)

	doPublish, err := askPublishTx("redeem")
	if err != nil {
		return err
	}
	if doPublish {
		txHash, err := xzc.Publish(*testnetFlag, rpcinfo, result.RedeemTx)
		if err != nil {
			return err
		}
		fmt.Printf("Published %s transaction (%s)\n", "redeem", txHash)
	}

	return nil
}

func refund(args []string) error {
	var rpcinfo libs.RPCInfo
	rpcinfo.HostPort = *connectFlag
	rpcinfo.User = *rpcuserFlag
	rpcinfo.Pass = *rpcpassFlag
	rpcinfo.WalletPass = *walletPass

	err := xzc.PingRPC(*testnetFlag, rpcinfo)
	if err != nil {
		return fmt.Errorf("Ping RPC: error: %v", err)
	}

	var params libs.RefundParams
	params.Contract = args[1]
	params.ContractTx = args[2]

	var result *libs.RefundResult
	result, err = xzc.Refund(*testnetFlag, rpcinfo, params)
	if err != nil {
		return fmt.Errorf("Refund: %v", err)
	}

	fmt.Printf("Refund fee: %d (%0.8f XZC/kB)\n\n", result.RefundFee, result.RefundFeePerKb)
	fmt.Printf("Refund transaction (%s):\n", result.RefundTxHash)
	fmt.Printf("%s\n\n", result.RefundTx)

	doPublish, err := askPublishTx("refund")
	if err != nil {
		return err
	}
	if doPublish {
		txHash, err := xzc.Publish(*testnetFlag, rpcinfo, result.RefundTx)
		if err != nil {
			return err
		}
		fmt.Printf("Published %s transaction (%s)\n", "refund", txHash)
	}

	return nil
}

func extractSecret(args []string) error {
	secret, err := xzc.ExtractSecret(args[1], args[2])
	if err != nil {
		return err
	}

	fmt.Printf("Contract shared secret: %s\n", secret)

	return nil
}

func auditContract(args []string) error {
	var params libs.AuditParams
	params.Contract = args[1]
	params.ContractTx = args[2]

	var result *libs.AuditResult
	result, err := xzc.AuditContract(*testnetFlag, params)
	if err != nil {
		return err
	}

	fmt.Printf("Contract address:        %s\n", result.ContractAddress)
	fmt.Printf("Contract value:          %v\n", xzc.Amount(result.ContractAmount))
	fmt.Printf("Recipient address:       %s\n", result.ContractRecipientAddress)
	fmt.Printf("Author's refund address: %s\n\n", result.ContractRefundAddress)

	fmt.Printf("Secret hash: %s\n\n", result.ContractSecretHash)

	locktime := result.ContractRefundLocktime
	if locktime >= int64(txscript.LockTimeThreshold) {
		t := time.Unix(locktime, 0)
		fmt.Printf("Locktime: %v\n", t.UTC())
		reachedAt := time.Until(t).Truncate(time.Second)
		if reachedAt > 0 {
			fmt.Printf("Locktime reached in %v\n", reachedAt)
		} else {
			fmt.Printf("Contract refund time lock has expired\n")
		}
	} else {
		fmt.Printf("Locktime: block %v\n", locktime)
	}

	return nil
}

func getTx(args []string) error {
	var rpcinfo libs.RPCInfo
	rpcinfo.HostPort = *connectFlag
	rpcinfo.User = *rpcuserFlag
	rpcinfo.Pass = *rpcpassFlag
	rpcinfo.WalletPass = *walletPass

	txid := args[1]

	result, err := xzc.GetTx(*testnetFlag, rpcinfo, txid)
	if err != nil {
		return fmt.Errorf("getTx: %v", err)
	}

	fmt.Printf("Confirmations: %d\n", result.Confirmations)
	blockHash := result.Blockhash
	if blockHash == "" {
		blockHash = "Unknown"
	}
	fmt.Printf("Block hash:    %s\n", blockHash)
	return nil
}

func newAddress(args []string) error {
	var rpcinfo libs.RPCInfo
	rpcinfo.HostPort = *connectFlag
	rpcinfo.User = *rpcuserFlag
	rpcinfo.Pass = *rpcpassFlag
	rpcinfo.WalletPass = *walletPass

	addr, err := xzc.GetNewAddress(*testnetFlag, rpcinfo)
	if err != nil {
		return fmt.Errorf("GetNewAddress: error: %v", err)
	}
	fmt.Printf("%s\n", addr)
	return nil
}

func askPublishTx(name string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Publish %s transaction? [y/N] ", name)
		answer, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		answer = strings.TrimSpace(strings.ToLower(answer))

		switch answer {
		case "y", "yes":
			return true, nil
		case "n", "no", "":
			return false, nil
		default:
			fmt.Println("please answer y or n")
			continue
		}
	}
}
