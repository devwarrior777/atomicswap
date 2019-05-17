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

// pingRPC is for the client to check if the server is running
func pingRPC(request *bnd.PingWalletRPCRequest) (*bnd.PingWalletRPCResponse, error) {
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

// newAddress gets a new address for the Coin & Network
func newAddress(request *bnd.NewAddressRequest) (*bnd.NewAddressResponse, error) {
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

// initiate gets contract data for the Coin & Network
func initiate(request *bnd.InitiateRequest) (*bnd.InitiateResponse, error) {
	conn, err := getClientConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := bnd.NewSwapLibClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, err := client.Initiate(ctx, request)
	if err != nil {
		return response, err
	}
	return response, nil
}

// participate gets contract data for the Coin & Network
func participate(request *bnd.ParticipateRequest) (*bnd.ParticipateResponse, error) {
	conn, err := getClientConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := bnd.NewSwapLibClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, err := client.Participate(ctx, request)
	if err != nil {
		return response, err
	}
	return response, nil
}

// redeem makes a redeem transaction for the contract, Coin & Network
func redeem(request *bnd.RedeemRequest) (*bnd.RedeemResponse, error) {
	conn, err := getClientConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := bnd.NewSwapLibClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, err := client.Redeem(ctx, request)
	if err != nil {
		return response, err
	}
	return response, nil
}

// refund makes a refund transaction for the contract, Coin & Network
func refund(request *bnd.RefundRequest) (*bnd.RefundResponse, error) {
	conn, err := getClientConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := bnd.NewSwapLibClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, err := client.Refund(ctx, request)
	if err != nil {
		return response, err
	}
	return response, nil
}

// extractSecret gets secret embedded in the initiator redeem tx scriptsig
func extractSecret(request *bnd.ExtractSecretRequest) (*bnd.ExtractSecretResponse, error) {
	conn, err := getClientConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := bnd.NewSwapLibClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, err := client.ExtractSecret(ctx, request)
	if err != nil {
		return response, err
	}
	return response, nil
}

// audit gets info embedded in the contract
func audit(request *bnd.AuditRequest) (*bnd.AuditResponse, error) {
	conn, err := getClientConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := bnd.NewSwapLibClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, err := client.Audit(ctx, request)
	if err != nil {
		return response, err
	}
	return response, nil
}

// publish broadcasts a transaction for the Coin & Network
func publish(request *bnd.PublishRequest) (*bnd.PublishResponse, error) {
	conn, err := getClientConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := bnd.NewSwapLibClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, err := client.Publish(ctx, request)
	if err != nil {
		return response, err
	}
	return response, nil
}

// gettx get info for a wallet's txid for the Coin & Network
func gettx(request *bnd.GetTxRequest) (*bnd.GetTxResponse, error) {
	conn, err := getClientConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := bnd.NewSwapLibClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, err := client.GetTx(ctx, request)
	if err != nil {
		return response, err
	}
	return response, nil
}

//...
