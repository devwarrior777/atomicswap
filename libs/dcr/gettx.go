package dcr

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/devwarrior777/atomicswap/libs"
)

const hexstr32 = 32 * 2

func getTx(testnet bool, rpcinfo libs.RPCInfo, txid string) (int32, string, error) {
	if len(txid) != hexstr32 {
		return 0, "", errors.New("txid: bad length")
	}
	txidBytes, err := hex.DecodeString(txid)
	if err != nil {
		return 0, "", err
	}
	wireTxHash := byteRev(txidBytes)

	wallet, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return 0, "", err
	}
	defer wallet.stopRPC()
	ctx := context.Background()

	var gtr *walletrpc.GetTransactionResponse
	gtr, err = wallet.client.GetTransaction(ctx, &walletrpc.GetTransactionRequest{
		TransactionHash: wireTxHash,
	})
	if err != nil {
		return 0, "", err
	}

	confirmations := gtr.Confirmations
	blockHash := hex.EncodeToString(byteRev(gtr.BlockHash))
	return confirmations, blockHash, nil
}
