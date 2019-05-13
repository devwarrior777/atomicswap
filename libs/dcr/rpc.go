// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package dcr

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/decred/dcrd/dcrutil"

	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/devwarrior777/atomicswap/libs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type wallet struct {
	conn   *grpc.ClientConn
	client walletrpc.WalletServiceClient
}

// starRPC - starts a new GRPC client for the network and address specified
//            along with the certs path, in RPCInfo
func startRPC(testnet bool, rpcinfo libs.RPCInfo) (*wallet, error) {
	hostport, err := getNormalizedAddress(testnet, rpcinfo.HostPort)
	if err != nil {
		return nil, fmt.Errorf("wallet server address: %v", err)
	}
	certPath := rpcinfo.Certs
	if certPath == "" {
		//default path
		certPath = filepath.Join(dcrutil.AppDataDir("dcrwallet", false), "rpc.cert")
	}
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		return nil, fmt.Errorf("open certificate: %v", err)
	}
	wallet := &wallet{}
	// get a connection to the server
	wallet.conn, err = grpc.Dial(hostport, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %v", err)
	}
	// get a client
	wallet.client = walletrpc.NewWalletServiceClient(wallet.conn)
	return wallet, err
}

// stopRPC closes the client connection
func (w *wallet) stopRPC() {
	w.conn.Close()
}

//////////////////////////////
// Miscellaneous GRPC funcs //
//////////////////////////////

func (w *wallet) ping() error {
	request := &walletrpc.PingRequest{}
	ctx := context.Background()
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	_, err := w.client.Ping(ctx, request)
	return err
}
