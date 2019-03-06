// Copyright (c) 2017/2019 The Decred developers
// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ltc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	rpc "github.com/ltcsuite/ltcd/rpcclient"
	"github.com/ltcsuite/ltcd/txscript"
	"github.com/ltcsuite/ltcd/wire"
	ltcutil "github.com/ltcsuite/ltcutil"
)

// startRPC - starts a new RPC client for the network and address specified
//            along with rpc user & rpc password, in RPCInfo
func startRPC(testnet bool, rpcinfo RPCInfo) (*rpc.Client, error) {
	hostport, err := getNormalizedAddress(testnet, rpcinfo.HostPort)
	if err != nil {
		return nil, fmt.Errorf("wallet server address: %v", err)
	}
	connConfig := &rpc.ConnConfig{
		Host:         hostport,
		User:         rpcinfo.User,
		Pass:         rpcinfo.Pass,
		DisableTLS:   true, //TODO: Should be a config option
		HTTPPostMode: true,
	}
	client, err := rpc.New(connConfig, nil)
	if err != nil {
		return client, fmt.Errorf("rpc connect: %v", err)
	}
	return client, err
}

// stopRPC - Explicit stop when not using defer()
func stopRPC(client *rpc.Client) {
	client.Shutdown()
	client.WaitForShutdown()
}

///////////////
// RPC funcs //
///////////////

// getNewAddress calls the getnewaddress JSON-RPC method.  It is
// implemented manually as the rpcclient implementation always passes the
// account parameter which was removed in Bitcoin Core 0.15.
func getNewAddress(testnet bool, rpcclient *rpc.Client) (ltcutil.Address, error) {
	chainParams := getChainParams(testnet)

	param0 := []json.RawMessage{[]byte(`""`)} // Deprecated but necessary in this position
	param1 := []json.RawMessage{[]byte(`"legacy"`)}
	params := append(param0, param1...)
	rawResp, err := rpcclient.RawRequest("getnewaddress", params)
	//
	if err != nil {
		return nil, err
	}
	var addrStr string
	err = json.Unmarshal(rawResp, &addrStr)
	if err != nil {
		return nil, err
	}
	addr, err := ltcutil.DecodeAddress(addrStr, chainParams)
	if err != nil {
		return nil, err
	}
	if !addr.IsForNet(chainParams) {
		return nil, fmt.Errorf("address %v is not intended for use on %v",
			addrStr, chainParams.Name)
	}
	if _, ok := addr.(*ltcutil.AddressPubKeyHash); !ok {
		return nil, fmt.Errorf("getnewaddress: address %v is not P2PKH",
			addr)
	}
	return addr, nil
}

// getRawChangeAddress calls the getrawchangeaddress JSON-RPC method.  It is
// implemented manually as the rpcclient implementation always passes the
// account parameter which was removed in Litecoin Core 0.15.
func getRawChangeAddress(testnet bool, rpcclient *rpc.Client) (ltcutil.Address, error) {
	chainParams := getChainParams(testnet)
	params := []json.RawMessage{[]byte(`"legacy"`)}
	rawResp, err := rpcclient.RawRequest("getrawchangeaddress", params)
	if err != nil {
		return nil, err
	}
	var addrStr string
	err = json.Unmarshal(rawResp, &addrStr)
	if err != nil {
		return nil, err
	}
	addr, err := ltcutil.DecodeAddress(addrStr, chainParams)
	if err != nil {
		return nil, err
	}
	if !addr.IsForNet(chainParams) {
		return nil, fmt.Errorf("address %v is not intended for use on %v",
			addrStr, chainParams.Name)
	}
	if _, ok := addr.(*ltcutil.AddressPubKeyHash); !ok {
		return nil, fmt.Errorf("getrawchangeaddress: address %v is not P2PKH",
			addr)
	}
	return addr, nil
}

// getFeePerKb queries the wallet for the transaction relay fee/kB to use and
// the minimum mempool relay fee.  It first tries to get the user-set fee in the
// wallet.  If unset, it attempts to find an estimate using estimatefee 6.  If
// both of these fail, it falls back to mempool relay fee policy.
func getFeePerKb(rpcclient *rpc.Client) (useFee, relayFee ltcutil.Amount, err error) {
	var netInfoResp struct {
		RelayFee float64 `json:"relayfee"`
	}
	var walletInfoResp struct {
		PayTxFee float64 `json:"paytxfee"`
	}
	var estimateResp struct {
		FeeRate float64 `json:"feerate"`
	}

	netInfoRawResp, err := rpcclient.RawRequest("getnetworkinfo", nil)
	if err == nil {
		err = json.Unmarshal(netInfoRawResp, &netInfoResp)
		if err != nil {
			return 0, 0, err
		}
	}
	walletInfoRawResp, err := rpcclient.RawRequest("getwalletinfo", nil)
	if err == nil {
		err = json.Unmarshal(walletInfoRawResp, &walletInfoResp)
		if err != nil {
			return 0, 0, err
		}
	}

	relayFee, err = ltcutil.NewAmount(netInfoResp.RelayFee)
	if err != nil {
		return 0, 0, err
	}
	payTxFee, err := ltcutil.NewAmount(walletInfoResp.PayTxFee)
	if err != nil {
		return 0, 0, err
	}

	// Use user-set wallet fee when set and not lower than the network relay
	// fee.
	if payTxFee != 0 {
		maxFee := payTxFee
		if relayFee > maxFee {
			maxFee = relayFee
		}
		return maxFee, relayFee, nil
	}

	params := []json.RawMessage{[]byte("6")}
	estimateRawResp, err := rpcclient.RawRequest("estimatesmartfee", params)
	if err != nil {
		return 0, 0, err
	}

	err = json.Unmarshal(estimateRawResp, &estimateResp)
	if err == nil && estimateResp.FeeRate > 0 {
		useFee, err = ltcutil.NewAmount(estimateResp.FeeRate)
		if relayFee > useFee {
			useFee = relayFee
		}
		return useFee, relayFee, err
	}

	fmt.Println("warning: falling back to mempool relay fee policy")
	return relayFee, relayFee, nil
}

// fundRawTransaction calls the fundrawtransaction JSON-RPC method.  It is
// implemented manually as client support is currently missing from the
// ltcd/rpcclient package.
func fundRawTransaction(rpcclient *rpc.Client, tx *wire.MsgTx, feePerKb ltcutil.Amount) (fundedTx *wire.MsgTx, fee ltcutil.Amount, err error) {
	var buf bytes.Buffer
	buf.Grow(tx.SerializeSize())
	tx.Serialize(&buf)
	param0, err := json.Marshal(hex.EncodeToString(buf.Bytes()))
	if err != nil {
		return nil, 0, err
	}
	param1, err := json.Marshal(struct {
		ChangeType string  `json:"change_type"`
		FeeRate    float64 `json:"feeRate"`
	}{
		ChangeType: "legacy",
		FeeRate:    feePerKb.ToBTC(),
	})
	if err != nil {
		return nil, 0, err
	}
	params := []json.RawMessage{param0, param1}
	rawResp, err := rpcclient.RawRequest("fundrawtransaction", params)
	if err != nil {
		return nil, 0, err
	}
	var resp struct {
		Hex       string  `json:"hex"`
		Fee       float64 `json:"fee"`
		ChangePos float64 `json:"changepos"`
	}
	err = json.Unmarshal(rawResp, &resp)
	if err != nil {
		return nil, 0, err
	}
	fundedTxBytes, err := hex.DecodeString(resp.Hex)
	if err != nil {
		return nil, 0, err
	}
	fundedTx = &wire.MsgTx{}
	err = fundedTx.Deserialize(bytes.NewReader(fundedTxBytes))
	if err != nil {
		return nil, 0, err
	}
	feeAmount, err := ltcutil.NewAmount(resp.Fee)
	if err != nil {
		return nil, 0, err
	}
	return fundedTx, feeAmount, nil
}

// createSig creates and returns the serialized raw signature and compressed
// pubkey for a transaction input signature.  Due to limitations of the Litecoin
// Core RPC API, this requires dumping a private key and signing in the client,
// rather than letting the wallet sign.
func createSig(testnet bool, tx *wire.MsgTx, idx int, pkScript []byte, addr ltcutil.Address,
	rpcclient *rpc.Client) (sig, pubkey []byte, err error) {

	wif, err := rpcclient.DumpPrivKey(addr)
	if err != nil {
		return nil, nil, err
	}
	sig, err = txscript.RawTxInSignature(tx, idx, pkScript, txscript.SigHashAll, wif.PrivKey)
	if err != nil {
		return nil, nil, err
	}
	return sig, wif.PrivKey.PubKey().SerializeCompressed(), nil
}

func sendRawTransaction(testnet bool, rpcclient *rpc.Client, tx *wire.MsgTx) (*chainhash.Hash, error) {
	txHash, err := rpcclient.SendRawTransaction(tx, false)
	if err != nil {
		return nil, fmt.Errorf("sendrawtransaction: %v", err)
	}
	return txHash, nil
}
