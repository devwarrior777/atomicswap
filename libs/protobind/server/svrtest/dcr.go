package svrtest

import (
	"errors"
	"fmt"

	"github.com/devwarrior777/atomicswap/libs"
	bnd "github.com/devwarrior777/atomicswap/libs/protobind"
	"google.golang.org/grpc/status"
)

func testDCR(testnet bool) error {
	// Store and re-use:
	//  - the address from NewAddress
	//  - contract and contract-tx from Initiate
	//  - generated secret hash for initiate, participate
	// We are testing the server here!
	var address string
	var contract string
	var contractTx string
	var redeemTx string
	var secret string
	var secretHash string

	// ping wallet
	pingreq := dcrPingWalletRPCRequest
	if testnet {
		pingreq = dcrTestnetPingWalletRPCRequest
	}
	ping, err := pingRPC(&pingreq)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf("status: %d - %v - %v", s.Code(), s.Code(), s.Message())
	}
	if ping.Errorno != bnd.ERRNO_OK {
		return fmt.Errorf("%v %s", ping.Errorno, ping.Errstr)
	}
	fmt.Println("Ping success")

	// new address
	newaddressreq := dcrNewAddressRequest
	if testnet {
		newaddressreq = dcrTestnetNewAddressRequest
	}
	newaddress, err := newAddress(&newaddressreq)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf("status: %d - %v - %v", s.Code(), s.Code(), s.Message())
	}
	if newaddress.Errorno != bnd.ERRNO_OK {
		return fmt.Errorf("%v %s", newaddress.Errorno, newaddress.Errstr)
	}
	address = newaddress.Address
	fmt.Printf("New address: %s\n", address)

	// initiate
	secret = libs.GetRand32()
	secretHash, err = libs.Hash256(secret)
	initiatereq := dcrInitiateRequest
	if testnet {
		initiatereq = dcrTestnetInitiateRequest
	}
	initiatereq.Secrethash = secretHash
	initiatereq.PartAddress = address
	initiate, err := initiate(&initiatereq)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf("status: %d - %v - %v", s.Code(), s.Code(), s.Message())
	}
	if initiate.Errorno != bnd.ERRNO_OK {
		return fmt.Errorf("%v %s", initiate.Errorno, initiate.Errstr)
	}
	contract = initiate.Contract
	contractTx = initiate.ContractTx
	if len(contract) < 64 || len(contractTx) < 64 {
		return errors.New("invalid contract/contract-tx length(s)")
	}
	fmt.Printf("Initiate contract:             %s...\n", contract[:64])
	fmt.Printf("Initiate contract tx:          %s...\n", contractTx[:64])
	fmt.Printf("Initiate P2SH address:         %s\n", initiate.ContractP2Sh)
	fmt.Printf("Initiate contract tx hash:     %s\n", initiate.ContractTxHash)
	fmt.Printf("Initiate fee:                  %d\n", initiate.Fee)
	fmt.Printf("Initiate fee rate:             %0.08f/kb\n", initiate.Feerate)
	fmt.Printf("Initiate refund locktime:      %d\n", initiate.Locktime)

	// participate
	participatereq := dcrParticipateRequest
	if testnet {
		participatereq = dcrTestnetParticipateRequest
	}
	participatereq.Secrethash = secretHash
	participatereq.InitAddress = address
	participate, err := participate(&participatereq)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf("status: %d - %v - %v", s.Code(), s.Code(), s.Message())
	}
	if participate.Errorno != bnd.ERRNO_OK {
		return fmt.Errorf("%v %s", participate.Errorno, participate.Errstr)
	}
	if len(participate.Contract) < 64 || len(participate.ContractTx) < 64 {
		return errors.New("invalid contract/contract-tx length(s)")
	}
	fmt.Printf("Participate contract:          %s...\n", participate.Contract[:64])
	fmt.Printf("Participate contract tx:       %s...\n", participate.ContractTx[:64])
	fmt.Printf("Participate P2SH address:      %s\n", participate.ContractP2Sh)
	fmt.Printf("Participate contract tx hash:  %s\n", participate.ContractTxHash)
	fmt.Printf("Participate fee:               %d\n", participate.Fee)
	fmt.Printf("Participate fee rate:          %0.08f/kb\n", participate.Feerate)
	fmt.Printf("Participate refund locktime:   %d\n", participate.Locktime)

	// audit
	auditreq := dcrAuditRequest
	if testnet {
		auditreq = dcrTestnetAuditRequest
	}
	auditreq.Contract = contract
	auditreq.ContractTx = contractTx
	audit, err := audit(&auditreq)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf("status: %d - %v - %v", s.Code(), s.Code(), s.Message())
	}
	if audit.Errorno == bnd.ERRNO_OK {
		// if here it is an error in the lib
		fmt.Printf("Audit contract amount:         %d\n", audit.ContractAmount)
		fmt.Printf("Audit contract address:        %s\n", audit.ContractAddress)
		fmt.Printf("Audit contract secret hash:    %s\n", audit.ContractSecrethash)
		fmt.Printf("Audit recipient address:       %s\n", audit.RecipientAddress)
		fmt.Printf("Audit refund address:          %s\n", audit.RefundAddress)
		fmt.Printf("Audit refund locktime:         %d\n", audit.RefundLocktime)
		return errors.New("got audit from an an unsupported lib function")
	}
	fmt.Printf("Expected error getting audit: %v %s\n", audit.Errorno, audit.Errstr)

	// redeem
	redeemreq := dcrRedeemRequest
	if testnet {
		redeemreq = dcrTestnetRedeemRequest
	}
	redeemreq.Secret = secret
	redeemreq.Contract = contract
	redeemreq.ContractTx = contractTx
	redeem, err := redeem(&redeemreq)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf("status: %d - %v - %v", s.Code(), s.Code(), s.Message())
	}
	if redeem.Errorno != bnd.ERRNO_OK {
		return fmt.Errorf("%v %s", redeem.Errorno, redeem.Errstr)
	}
	redeemTx = redeem.RedeemTx
	if len(redeemTx) < 64 {
		return errors.New("invalid contract/contract-tx length(s)")
	}
	fmt.Printf("Redeem contract:               %s...\n", redeemTx[:64])
	fmt.Printf("Redeem contract tx:            %s...\n", redeem.RedeemTxHash)
	fmt.Printf("Redeem fee:                    %d\n", redeem.Fee)
	fmt.Printf("Redeem fee rate:               %0.08f/kb\n", redeem.Feerate)

	// extractSecret
	extractsecretreq := dcrExtractSecretRequest
	if testnet {
		extractsecretreq = dcrTestnetExtractSecretRequest
	}
	extractsecretreq.CpRedemptionTx = redeemTx
	extractsecretreq.Secrethash = secretHash
	extract, err := extractSecret(&extractsecretreq)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf("status: %d - %v - %v", s.Code(), s.Code(), s.Message())
	}
	if extract.Errorno == bnd.ERRNO_OK {
		// if here it is an error in the lib
		fmt.Printf("ExtractSecret secret:          %s\n", extract.Secret)
		return errors.New("extractSecret from an an unsupported lib function")
	}
	fmt.Printf("Expected error getting audit: %v %s\n", extract.Errorno, extract.Errstr)

	// refund
	refundreq := dcrRefundRequest
	if testnet {
		refundreq = dcrTestnetRefundRequest
	}
	refundreq.Contract = contract
	refundreq.ContractTx = contractTx
	refund, err := refund(&refundreq)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf("status: %d - %v - %v", s.Code(), s.Code(), s.Message())
	}
	if refund.Errorno != bnd.ERRNO_OK {
		return fmt.Errorf("%v %s", refund.Errorno, refund.Errstr)
	}
	if len(refund.RefundTx) < 64 {
		return errors.New("invalid contract/contract-tx length(s)")
	}
	fmt.Printf("Refund tx:                     %s...\n", refund.RefundTx[:64])
	fmt.Printf("Refund tx hash:                %s...\n", refund.RefundTxHash)
	fmt.Printf("Refund fee:                    %d\n", refund.Fee)
	fmt.Printf("Refund fee rate:               %0.08f/kb\n", refund.Feerate)

	// publish
	//
	// This a negative test since we do not want to boadcast a transaction to
	// the network.
	//
	// It tests that the test client can reach and call the wallet node publish
	// function through the gRPC server
	//
	publishreq := dcrPublishRequest
	if testnet {
		publishreq = dcrTestnetPublishRequest
	}
	publishreq.Tx = "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
	publish, err := publish(&publishreq)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf("status: %d - %v - %v", s.Code(), s.Code(), s.Message())
	}
	if publish.Errorno == bnd.ERRNO_OK {
		// if here it is an error in the lib
		fmt.Printf("Publish contract:               %s...\n", publish.TxHash)
		return errors.New("published invalid transaction")
	}
	fmt.Printf("Expected error publishing invalid transaction: %v %s\n", publish.Errorno, publish.Errstr)

	// gettx
	//
	// This a negative test again
	//
	// It tests that the test client can reach and call the wallet node gettx
	// function through the gRPC server
	//
	gettxreq := dcrGetTxRequest
	if testnet {
		gettxreq = dcrTestnetGetTxRequest
	}
	gettxreq.Txid = "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
	gettx, err := gettx(&gettxreq)
	if err != nil {
		s := status.Convert(err)
		return fmt.Errorf("status: %d - %v - %v", s.Code(), s.Code(), s.Message())
	}
	if gettx.Errorno == bnd.ERRNO_OK {
		// if here it is an error in the lib
		fmt.Printf("GetTx Blockhash:               %s...\n", gettx.Blockhash)
		return errors.New("got info from an invalid txid")
	}
	fmt.Printf("Expected error getting info from an invalid txid: %v %s\n", gettx.Errorno, gettx.Errstr)

	return nil
}
