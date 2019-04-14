package xzc

import (
	"github.com/devwarrior777/atomicswap/libs"
)

func getTx(testnet bool, rpcinfo libs.RPCInfo, txid string) (*libs.GetTxResult, error) {
	rpcclient, err := startRPC(testnet, rpcinfo)
	if err != nil {
		return nil, err
	}
	defer func() {
		rpcclient.Shutdown()
		rpcclient.WaitForShutdown()
	}()

	result, err := getTransaction(rpcclient, txid)
	if err != nil {
		return nil, err
	}

	return result, nil
}
