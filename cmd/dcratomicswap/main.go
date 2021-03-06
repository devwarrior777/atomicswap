// Copyright (c) 2017 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	"github.com/devwarrior777/atomicswap/libs"
	"github.com/devwarrior777/atomicswap/libs/dcr"
)

var (
	flagset     = flag.NewFlagSet("", flag.ExitOnError)
	connectFlag = flagset.String("s", "localhost", "host[:port] of dcrwallet gRPC server")
	certFlag    = flagset.String("c", filepath.Join(dcrutil.AppDataDir("dcrwallet", false), "rpc.cert"), "dcrwallet RPC certificate path")
	testnetFlag = flagset.Bool("testnet", false, "use testnet network")
	walletPass  = flagset.String("wpass", "", "wallet passphrase")
)

// There are two directions that the atomic swap can be performed, as the
// initiator can be on either chain.  This tool only deals with creating the
// Decred transactions for these swaps.  A second tool should be used for the
// transaction on the other chain.  Any chain can be used so long as it supports
// OP_SHA256 and OP_CHECKLOCKTIMEVERIFY.
//
// Example scenerios using bitcoin as the second chain:
//
// Scenerio 1:
//   cp1 initiates (dcr)
//   cp2 participates with cp1 H(S) (btc)
//   cp1 redeems btc revealing S
//     - must verify H(S) in contract is hash of known secret
//   cp2 redeems dcr with S
//
// Scenerio 2:
//   cp1 initiates (btc)
//   cp2 participates with cp1 H(S) (dcr)
//   cp1 redeems dcr revealing S
//     - must verify H(S) in contract is hash of known secret
//   cp2 redeems btc with S

func init() {
	flagset.Usage = func() {
		fmt.Println("Usage: dcratomicswap [flags] cmd [cmd args]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  initiate <participant address> <amount>")
		fmt.Println("  participate <initiator address> <amount> <secret hash>")
		fmt.Println("  redeem <contract> <contract transaction> <secret>")
		fmt.Println("  refund <contract> <contract transaction>")
		fmt.Println("  extractsecret <redemption transaction> <secret hash> [NOT IMPLEMENTED]")
		fmt.Println("  auditcontract <contract> <contract transaction>      [NOT IMPLEMENTED]")
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

	amount, err := dcr.NewAmount(amountF64)
	if err != nil {
		return err
	}

	var rpcinfo libs.RPCInfo
	rpcinfo.HostPort = *connectFlag
	rpcinfo.WalletPass = *walletPass
	rpcinfo.Certs = *certFlag

	err = dcr.PingRPC(*testnetFlag, rpcinfo)
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
	result, err = dcr.Initiate(*testnetFlag, rpcinfo, params)
	if err != nil {
		return fmt.Errorf("Initiate: %v", err)
	}

	fmt.Printf("Secret:      %s\n", secret)
	fmt.Printf("Secret hash: %s\n\n", secretHash)
	fmt.Printf("Contract fee: %d (%0.8f DCR/kB)\n", result.ContractFee, result.ContractFeePerKb)
	fmt.Printf("Contract (%s):\n", result.ContractP2SH)
	fmt.Printf("%s\n\n", result.Contract)
	fmt.Printf("Contract transaction (%s):\n", result.ContractTxHash)
	fmt.Printf("%s\n\n", result.ContractTx)

	doPublish, err := askPublishTx("contract")
	if err != nil {
		return err
	}
	if doPublish {
		txHash, err := dcr.Publish(*testnetFlag, rpcinfo, result.ContractTx)
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
	amount, err := dcr.NewAmount(amountF64)
	if err != nil {
		return err
	}

	var rpcinfo libs.RPCInfo
	rpcinfo.HostPort = *connectFlag
	rpcinfo.WalletPass = *walletPass
	rpcinfo.Certs = *certFlag

	err = dcr.PingRPC(*testnetFlag, rpcinfo)
	if err != nil {
		return fmt.Errorf("Ping RPC: error: %v", err)
	}

	var params libs.ParticipateParams
	params.SecretHash = args[3]
	params.CP1Addr = args[1]
	params.CP1Amount = int64(amount)

	var result *libs.ParticipateResult
	result, err = dcr.Participate(*testnetFlag, rpcinfo, params)
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
		txHash, err := dcr.Publish(*testnetFlag, rpcinfo, result.ContractTx)
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
	rpcinfo.Certs = *certFlag
	rpcinfo.WalletPass = *walletPass

	err := dcr.PingRPC(*testnetFlag, rpcinfo)
	if err != nil {
		return fmt.Errorf("Ping RPC: error: %v", err)
	}

	var params libs.RedeemParams
	params.Contract = args[1]
	params.ContractTx = args[2]
	params.Secret = args[3]

	var result *libs.RedeemResult
	result, err = dcr.Redeem(*testnetFlag, rpcinfo, params)
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
		txHash, err := dcr.Publish(*testnetFlag, rpcinfo, result.RedeemTx)
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
	rpcinfo.Certs = *certFlag
	rpcinfo.WalletPass = *walletPass

	err := dcr.PingRPC(*testnetFlag, rpcinfo)
	if err != nil {
		return fmt.Errorf("Ping RPC: error: %v", err)
	}

	var params libs.RefundParams
	params.Contract = args[1]
	params.ContractTx = args[2]

	var result *libs.RefundResult
	result, err = dcr.Refund(*testnetFlag, rpcinfo, params)
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
		txHash, err := dcr.Publish(*testnetFlag, rpcinfo, result.RefundTx)
		if err != nil {
			return err
		}
		fmt.Printf("Published %s transaction (%s)\n", "refund", txHash)
	}

	return nil
}

func extractSecret(args []string) error {
	fmt.Println("Not implemented")
	return nil
}

func auditContract(args []string) error {
	fmt.Println("Not implemented")
	return nil
}

func getTx(args []string) error {
	var rpcinfo libs.RPCInfo
	rpcinfo.HostPort = *connectFlag
	rpcinfo.Certs = *certFlag
	rpcinfo.WalletPass = *walletPass

	txid := args[1]

	result, err := dcr.GetTx(*testnetFlag, rpcinfo, txid)
	if err != nil {
		return fmt.Errorf("getTx: %v", err)
	}
	fmt.Printf("Confirmations: %d\n", result.Confirmations)
	fmt.Printf("Block hash:    %s\n", result.Blockhash)
	return nil
}

func newAddress(args []string) error {
	var rpcinfo libs.RPCInfo
	rpcinfo.HostPort = *connectFlag
	rpcinfo.Certs = *certFlag
	rpcinfo.WalletPass = *walletPass

	addr, err := dcr.GetNewAddress(*testnetFlag, rpcinfo)
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
