package svrtest

import (
	"context"
	"fmt"
	"log"

	bnd "github.com/devwarrior777/atomicswap/libs/protobind"
	"github.com/devwarrior777/atomicswap/libs/protobind/server/svrcfg"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	useTLS             = svrcfg.Config.UseTLS
	certPath           = svrcfg.Config.CertPath
	serverAddr         = svrcfg.Config.ServerAddr
	serverPort         = svrcfg.Config.ServerPort
	serverHostOverride = svrcfg.Config.HostOverride
)

// getClientConnection gets a connection to the swap session server
func getClientConnection() (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if useTLS {
		creds, err := credentials.NewClientTLSFromFile(certPath, serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		log.Println("Warning: No TLS")
		opts = append(opts, grpc.WithInsecure())
	}
	serverHostPort := fmt.Sprintf("%s:%d", serverAddr, serverPort)
	conn, err := grpc.Dial(serverHostPort, opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

/////////
// API //
/////////

// PingRPC is for the client to check if the server is running
func PingRPC(request *bnd.PingWalletRPCRequest) (*bnd.PingWalletRPCResponse, error) {
	conn, err := getClientConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := bnd.NewSwapLibClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, err := client.PingWalletRPC(ctx, request)
	if err != nil {
		return response, err
	}
	return response, nil
}

// NewAddress gets a new address for the Coin & Network
func NewAddress(request *bnd.NewAddressRequest) (*bnd.NewAddressResponse, error) {
	conn, err := getClientConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := bnd.NewSwapLibClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, err := client.NewAddress(ctx, request)
	if err != nil {
		return response, err
	}
	return response, nil
}

//...
