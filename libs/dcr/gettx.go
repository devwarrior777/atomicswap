package dcr

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/devwarrior777/atomicswap/libs"
)

const hexstr32 = 32 * 2

func getTx(testnet bool, rpcinfo libs.RPCInfo, txid string) (*libs.GetTxResult, error) {
	if len(txid) != hexstr32 {
		return nil, errors.New("txid: bad length")
	}
	txidBytes, err := hex.DecodeString(txid)
	if err != nil {
		return nil, err
	}
	wireTxHash := byteRev(txidBytes)

	wallet, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return nil, err
	}
	defer wallet.stopRPC()
	ctx := context.Background()

	var gtr *walletrpc.GetTransactionResponse
	gtr, err = wallet.client.GetTransaction(ctx, &walletrpc.GetTransactionRequest{
		TransactionHash: wireTxHash,
	})
	if err != nil {
		return nil, err
	}

	result := &libs.GetTxResult{}
	result.Confirmations = uint64(gtr.Confirmations)
	result.Blockhash = hex.EncodeToString(byteRev(gtr.BlockHash))

	return result, nil
}
