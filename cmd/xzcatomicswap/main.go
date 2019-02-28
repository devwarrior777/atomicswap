// Copyright (c) 2017/2019 The Decred developers
// Copyright (c) 2018/2019 The Zcoin developers
// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/devwarrior777/atomicswap/libs/xzc" // Use new libs/xzc pkg
	"github.com/zcoinofficial/xzcd/chaincfg"
	"github.com/zcoinofficial/xzcd/txscript"
	"github.com/zcoinofficial/xzcd/wire"
	xzcutil "github.com/zcoinofficial/xzcutil"
)

var (
	chainParams = &chaincfg.MainNetParams
)

var (
	flagset     = flag.NewFlagSet("", flag.ExitOnError)
	connectFlag = flagset.String("s", "localhost", "host[:port] of Zcoin Core wallet RPC server")
	rpcuserFlag = flagset.String("rpcuser", "", "username for wallet RPC authentication")
	rpcpassFlag = flagset.String("rpcpass", "", "password for wallet RPC authentication")
	testnetFlag = flagset.Bool("testnet", false, "use testnet network")
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
		fmt.Println()
		fmt.Println("Flags:")
		flagset.PrintDefaults()
	}
}

func main() {
	showUsage, err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if showUsage {
		flagset.Usage()
	}
	if err != nil || showUsage {
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

func run() (showUsage bool, err error) {
	flagset.Parse(os.Args[1:])
	args := flagset.Args()
	if len(args) == 0 {
		return true, nil
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
	default:
		return true, fmt.Errorf("unknown command %v", args[0])
	}
	nArgs := checkCmdArgLength(args[1:], cmdArgs)
	flagset.Parse(args[1+nArgs:])
	if nArgs < cmdArgs {
		return true, fmt.Errorf("%s: too few arguments", args[0])
	}
	if flagset.NArg() != 0 {
		return true, fmt.Errorf("unexpected argument: %s", flagset.Arg(0))
	}

	if *testnetFlag {
		chainParams = &chaincfg.TestNet3Params
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
	}

	return true, fmt.Errorf("unexpected argument: %s", flagset.Arg(0))
}

func initiate(args []string) (bool, error) {
	showUsage := true
	cp2Addr, err := xzcutil.DecodeAddress(args[1], chainParams)
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode participant address: %v", err)
	}
	if !cp2Addr.IsForNet(chainParams) {
		return showUsage, fmt.Errorf("participant address is not "+
			"intended for use on %v", chainParams.Name)
	}
	cp2AddrP2PKH, ok := cp2Addr.(*xzcutil.AddressPubKeyHash)
	if !ok {
		return showUsage, errors.New("participant address is not P2PKH")
	}

	amountF64, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode amount: %v", err)
	}

	showUsage = false

	amount, err := xzcutil.NewAmount(amountF64)
	if err != nil {
		return showUsage, err
	}

	var rpcinfo xzc.RPCInfo
	rpcinfo.HostPort = *connectFlag
	rpcinfo.User = *rpcuserFlag
	rpcinfo.Pass = *rpcpassFlag
	var params xzc.InitiateParams
	params.CP2AddrP2PKH = cp2AddrP2PKH
	params.CP2Amount = amount

	var result xzc.InitiateResult
	result, err = xzc.Initiate(*testnetFlag, rpcinfo, params)
	if err != nil {
		return showUsage, fmt.Errorf("Initiate: %v", err)
	}

	var refundParams xzc.RefundParams
	refundParams.Contract = result.Contract
	refundParams.ContractTx = &result.ContractTx

	var refundResult xzc.RefundResult
	refundResult, err = xzc.Refund(*testnetFlag, rpcinfo, refundParams)
	if err != nil {
		return showUsage, fmt.Errorf("Initiate: %v", err)
	}

	contractTxHash := result.ContractTx.TxHash()
	refundTxHash := refundResult.RefundTx.TxHash()

	fmt.Printf("Secret:      %x\n", result.Secret)
	fmt.Printf("Secret hash: %x\n\n", result.SecretHash)
	fmt.Printf("Contract fee: %v (%0.8f XZC/kB)\n", result.ContractFee, result.ContractFeePerKb)
	fmt.Printf("Refund fee:   %v (%0.8f XZC/kB)\n\n", refundResult.RefundFee, refundResult.RefundFeePerKb)
	fmt.Printf("Contract (%v):\n", result.ContractP2SH)
	fmt.Printf("%x\n\n", result.Contract)
	var contractBuf bytes.Buffer
	contractBuf.Grow(result.ContractTx.SerializeSize())
	result.ContractTx.Serialize(&contractBuf)
	fmt.Printf("Contract transaction (%v):\n", contractTxHash)
	fmt.Printf("%x\n\n", contractBuf.Bytes())
	var refundBuf bytes.Buffer
	refundBuf.Grow(refundResult.RefundTx.SerializeSize())
	refundResult.RefundTx.Serialize(&refundBuf)
	fmt.Printf("Refund transaction (%v):\n", &refundTxHash)
	fmt.Printf("%x\n\n", refundBuf.Bytes())

	doPublish, err := askPublishTx("contract")
	if err != nil {
		return showUsage, err
	}
	if doPublish {
		txHash, err := xzc.Publish(*testnetFlag, rpcinfo, &result.ContractTx)
		if err != nil {
			return showUsage, err
		}
		fmt.Printf("Published %s transaction (%v)\n", "contract", txHash)
	}

	return showUsage, nil
}

func participate(args []string) (bool, error) {
	showUsage := true
	cp1Addr, err := xzcutil.DecodeAddress(args[1], chainParams)
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode initiator address: %v", err)
	}
	if !cp1Addr.IsForNet(chainParams) {
		return showUsage, fmt.Errorf("initiator address is not "+
			"intended for use on %v", chainParams.Name)
	}
	cp1AddrP2PKH, ok := cp1Addr.(*xzcutil.AddressPubKeyHash)
	if !ok {
		return showUsage, errors.New("initiator address is not P2PKH")
	}

	amountF64, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode amount: %v", err)
	}
	amount, err := xzcutil.NewAmount(amountF64)
	if err != nil {
		return showUsage, err
	}

	secretHash, err := hex.DecodeString(args[3])
	if err != nil {
		return showUsage, errors.New("secret hash must be hex encoded")
	}
	if len(secretHash) != sha256.Size {
		return showUsage, errors.New("secret hash has wrong size")
	}

	showUsage = false

	var rpcinfo xzc.RPCInfo
	rpcinfo.HostPort = *connectFlag
	rpcinfo.User = *rpcuserFlag
	rpcinfo.Pass = *rpcpassFlag
	var params xzc.ParticipateParams
	params.SecretHash = secretHash
	params.CP1AddrP2PKH = cp1AddrP2PKH
	params.CP1Amount = amount

	var result xzc.ParticipateResult
	result, err = xzc.Participate(*testnetFlag, rpcinfo, params)
	if err != nil {
		return showUsage, fmt.Errorf("Participate: %v", err)
	}

	var refundParams xzc.RefundParams
	refundParams.Contract = result.Contract
	refundParams.ContractTx = &result.ContractTx

	var refundResult xzc.RefundResult
	refundResult, err = xzc.Refund(*testnetFlag, rpcinfo, refundParams)
	if err != nil {
		return showUsage, fmt.Errorf("Refund: %v", err)
	}

	contractTxHash := result.ContractTx.TxHash()
	refundTxHash := refundResult.RefundTx.TxHash()

	fmt.Printf("Contract fee: %v (%0.8f XZC/kB)\n", result.ContractFee, result.ContractFeePerKb)
	fmt.Printf("Refund fee:   %v (%0.8f XZC/kB)\n\n", refundResult.RefundFee, refundResult.RefundFeePerKb)
	fmt.Printf("Contract (%v):\n", result.ContractP2SH)
	fmt.Printf("%x\n\n", result.Contract)
	var contractBuf bytes.Buffer
	contractBuf.Grow(result.ContractTx.SerializeSize())
	result.ContractTx.Serialize(&contractBuf)
	fmt.Printf("Contract transaction (%v):\n", contractTxHash)
	fmt.Printf("%x\n\n", contractBuf.Bytes())

	var refundBuf bytes.Buffer
	refundBuf.Grow(refundResult.RefundTx.SerializeSize())
	refundResult.RefundTx.Serialize(&refundBuf)
	fmt.Printf("Refund transaction (%v):\n", &refundTxHash)
	fmt.Printf("%x\n\n", refundBuf.Bytes())

	doPublish, err := askPublishTx("contract")
	if err != nil {
		return showUsage, err
	}
	if doPublish {
		txHash, err := xzc.Publish(*testnetFlag, rpcinfo, &result.ContractTx)
		if err != nil {
			return showUsage, err
		}
		fmt.Printf("Published %s transaction (%v)\n", "contract", txHash)
	}

	return showUsage, nil
}

func redeem(args []string) (bool, error) {
	showUsage := true
	contract, err := hex.DecodeString(args[1])
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode contract: %v", err)
	}

	contractTxBytes, err := hex.DecodeString(args[2])
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode contract transaction: %v", err)
	}
	var contractTx wire.MsgTx
	err = contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode contract transaction: %v", err)
	}

	secret, err := hex.DecodeString(args[3])
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode secret: %v", err)
	}

	showUsage = false

	var rpcinfo xzc.RPCInfo
	rpcinfo.HostPort = *connectFlag
	rpcinfo.User = *rpcuserFlag
	rpcinfo.Pass = *rpcpassFlag
	var params xzc.RedeemParams
	params.Secret = secret
	params.Contract = contract
	params.ContractTx = &contractTx

	var result xzc.RedeemResult
	result, err = xzc.Redeem(*testnetFlag, rpcinfo, params)
	if err != nil {
		return showUsage, fmt.Errorf("Redeem: %v", err)
	}

	redeemTx := result.RedeemTx
	redeemTxHash := redeemTx.TxHash()

	var buf bytes.Buffer
	buf.Grow(redeemTx.SerializeSize())
	redeemTx.Serialize(&buf)
	fmt.Printf("Redeem fee: %v (%0.8f XZC/kB)\n\n", result.RedeemFee, result.RedeemFeePerKb)
	fmt.Printf("Redeem transaction (%v):\n", redeemTxHash)
	fmt.Printf("%x\n\n", buf.Bytes())

	doPublish, err := askPublishTx("redeem")
	if err != nil {
		return showUsage, err
	}
	if doPublish {
		txHash, err := xzc.Publish(*testnetFlag, rpcinfo, &redeemTx)
		if err != nil {
			return false, err
		}
		fmt.Printf("Published %s transaction (%v)\n", "redeem", txHash)
	}

	return false, nil
}

func refund(args []string) (bool, error) {
	showUsage := true
	contract, err := hex.DecodeString(args[1])
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode contract: %v", err)
	}

	contractTxBytes, err := hex.DecodeString(args[2])
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode contract transaction: %v", err)
	}
	var contractTx wire.MsgTx
	err = contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode contract transaction: %v", err)
	}

	showUsage = false

	var rpcinfo xzc.RPCInfo
	rpcinfo.HostPort = *connectFlag
	rpcinfo.User = *rpcuserFlag
	rpcinfo.Pass = *rpcpassFlag
	var params xzc.RefundParams
	params.Contract = contract
	params.ContractTx = &contractTx

	var result xzc.RefundResult
	result, err = xzc.Refund(*testnetFlag, rpcinfo, params)
	if err != nil {
		return showUsage, fmt.Errorf("Refund: %v", err)
	}

	refundTx := result.RefundTx
	refundTxHash := refundTx.TxHash()

	var buf bytes.Buffer
	buf.Grow(refundTx.SerializeSize())
	refundTx.Serialize(&buf)
	fmt.Printf("Refund fee: %v (%0.8f XZC/kB)\n\n", result.RefundFee, result.RefundFeePerKb)
	fmt.Printf("Refund transaction (%v):\n", refundTxHash)
	fmt.Printf("%x\n\n", buf.Bytes())

	doPublish, err := askPublishTx("refund")
	if err != nil {
		return showUsage, err
	}
	if doPublish {
		txHash, err := xzc.Publish(*testnetFlag, rpcinfo, &refundTx)
		if err != nil {
			return showUsage, err
		}
		fmt.Printf("Published %s transaction (%v)\n", "refund", txHash)
	}

	return showUsage, nil
}

func extractSecret(args []string) (bool, error) {
	showUsage := true
	redemptionTxBytes, err := hex.DecodeString(args[1])
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode redemption transaction: %v", err)
	}
	var redemptionTx wire.MsgTx
	err = redemptionTx.Deserialize(bytes.NewReader(redemptionTxBytes))
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode redemption transaction: %v", err)
	}

	secretHash, err := hex.DecodeString(args[2])
	if err != nil {
		return showUsage, errors.New("secret hash must be hex encoded")
	}
	if len(secretHash) != sha256.Size {
		return showUsage, errors.New("secret hash has wrong size")
	}

	showUsage = false

	secret, err := xzc.ExtractSecret(&redemptionTx, secretHash)
	if err != nil {
		return showUsage, err
	}

	fmt.Printf("Contract shared secret: %x\n", secret)

	return showUsage, nil
}

func auditContract(args []string) (bool, error) {
	showUsage := true
	contract, err := hex.DecodeString(args[1])
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode contract: %v", err)
	}

	contractTxBytes, err := hex.DecodeString(args[2])
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode contract transaction: %v", err)
	}
	var contractTx wire.MsgTx
	err = contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		return showUsage, fmt.Errorf("failed to decode contract transaction: %v", err)
	}

	showUsage = false

	var params xzc.AuditParams
	params.Contract = contract
	params.ContractTx = &contractTx

	var result xzc.AuditResult
	result, err = xzc.AuditContract(*testnetFlag, params)
	if err != nil {
		return false, err
	}

	fmt.Printf("Contract address:        %v\n", result.ContractAddress.EncodeAddress())
	fmt.Printf("Contract value:          %v\n", result.ContractAmount)
	fmt.Printf("Recipient address:       %v\n", result.ContractRecipientAddress.EncodeAddress())
	fmt.Printf("Author's refund address: %v\n\n", result.ContractRefundAddress.EncodeAddress())

	fmt.Printf("Secret hash: %x\n\n", result.ContractSecretHash)

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

	return false, nil
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
