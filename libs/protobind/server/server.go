package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"github.com/devwarrior777/atomicswap/libs"
	bnd "github.com/devwarrior777/atomicswap/libs/protobind"
	"github.com/devwarrior777/atomicswap/libs/protobind/server/svrcfg"
	"github.com/devwarrior777/atomicswap/libs/protobind/server/wallets"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Configuration
var (
	pidFile     = svrcfg.Config.PidFile
	tls         = svrcfg.Config.UseTLS
	certPath    = svrcfg.Config.CertPath
	certKeyPath = svrcfg.Config.CertKeyPath
	serverAddr  = svrcfg.Config.ServerAddr
	serverPort  = svrcfg.Config.ServerPort
)

// gRPC server instance
var grpcServer *grpc.Server

// swapLibServer implements swapLibServer
type swapLibServer struct {
}

///////////////////////////////////
// Only meta response will error //
///////////////////////////////////

// PingWalletRPC pings the wallet node RPC client to establish if the node is running
func (s *swapLibServer) PingWalletRPC(ctx context.Context, request *bnd.PingWalletRPCRequest) (*bnd.PingWalletRPCResponse, error) {
	log.Printf("PingWalletRPC\n")
	response := &bnd.PingWalletRPCResponse{Errorno: bnd.ERRNO_OK}
	// get wallet
	rpcinfo := libs.RPCInfo{}
	rpcinfo.HostPort = request.Hostport
	rpcinfo.User = request.Rpcuser
	rpcinfo.Pass = request.Rpcpass
	rpcinfo.WalletPass = request.Wpass
	rpcinfo.Certs = request.Certs
	wallet, err := wallets.WalletForCoin(request.Testnet, rpcinfo, request.Coin)
	if err != nil {
		response.Errorno = bnd.ERRNO_LIBS
		response.Errstr = err.Error()
		return response, nil
	}
	// ping wallet
	err = wallet.PingRPC()
	if err != nil {
		response.Errorno = bnd.ERRNO_LIBS
		response.Errstr = err.Error()
		return response, nil
	}
	return response, nil
}

func (s *swapLibServer) NewAddress(ctx context.Context, request *bnd.NewAddressRequest) (*bnd.NewAddressResponse, error) {
	// response := newAddress(request)
	response := &bnd.NewAddressResponse{}
	log.Printf("NewAddress\n")
	return response, nil
}

func (s *swapLibServer) Initiate(ctx context.Context, request *bnd.InitiateRequest) (*bnd.InitiateResponse, error) {
	// response := initiate(request)
	response := &bnd.InitiateResponse{}
	log.Printf("Initiate\n")
	return response, nil
}

func (s *swapLibServer) Participate(ctx context.Context, request *bnd.ParticipateRequest) (*bnd.ParticipateResponse, error) {
	// response := participate(request)
	response := &bnd.ParticipateResponse{}
	log.Printf("Initiate\n")
	return response, nil
}

func (s *swapLibServer) Redeem(ctx context.Context, request *bnd.RedeemRequest) (*bnd.RedeemResponse, error) {
	// response := redeem(request)
	response := &bnd.RedeemResponse{}
	log.Printf("Redeem\n")
	return response, nil
}

func (s *swapLibServer) Refund(ctx context.Context, request *bnd.RefundRequest) (*bnd.RefundResponse, error) {
	// response := refund(request)
	response := &bnd.RefundResponse{}
	log.Printf("Refund\n")
	return response, nil
}

func (s *swapLibServer) Publish(ctx context.Context, request *bnd.PublishRequest) (*bnd.PublishResponse, error) {
	// response := publish(request)
	response := &bnd.PublishResponse{}
	log.Printf("Publish\n")
	return response, nil
}

func (s *swapLibServer) ExtractSecret(ctx context.Context, request *bnd.ExtractSecretRequest) (*bnd.ExtractSecretResponse, error) {
	// response := extractSecret(request)
	response := &bnd.ExtractSecretResponse{}
	log.Printf("ExtractSecret\n")
	return response, nil
}

func (s *swapLibServer) Audit(ctx context.Context, request *bnd.AuditRequest) (*bnd.AuditResponse, error) {
	// response := audit(request)
	response := &bnd.AuditResponse{}
	log.Printf("Audit\n")
	return response, nil
}

func (s *swapLibServer) GetTx(ctx context.Context, request *bnd.GetTxRequest) (*bnd.GetTxResponse, error) {
	// response := getTx(request)
	response := &bnd.GetTxResponse{}
	log.Printf("GetTx\n")
	return response, nil
}

//////////
// MAIN //
//////////

// newServer is the swapLibServer Constructor
func newServer() *swapLibServer {
	s := &swapLibServer{}
	return s
}

func main() {
	ensureUniqueServerProcess()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", serverPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Server listening on localhost:%v\n", serverPort)
	var opts []grpc.ServerOption
	if tls {
		creds, err := credentials.NewServerTLSFromFile(certPath, certKeyPath)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	} else {
		log.Println("Warning: No TLS")
	}
	// export process lock/pid file
	setPidFile()
	// Good to go
	startSignalHandler()
	grpcServer = grpc.NewServer(opts...)
	bnd.RegisterSwapLibServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

/////////////////////////
// One Server Instance //
/////////////////////////

func ensureUniqueServerProcess() {
	if runtime.GOOS == "windows" {
		log.Fatalln("This server does not run on Windows")
	}
	if checkPidfileExists() {
		log.Fatalln("server already running")
	}
}

func checkPidfileExists() bool {
	_, err := os.Stat(pidFile)
	if err == nil {
		return true
	}
	// log.Printf("checkPidfileExists: %v\n", err)
	return false
}

//////////////////
// Needed Setup //
//////////////////

// Allow shutdown process to discover and gracefully stop server
func setPidFile() {
	pid := strconv.FormatInt(int64(os.Getpid()), 10)
	f, err := os.Create(pidFile)
	if err != nil {
		log.Fatalf("cannot create pid file: %s\n", pidFile)
	}
	defer f.Close()
	f.WriteString(pid)
	f.Sync()
	log.Printf("server pid: %s\n", pid)
}

///////////////////////////////////////
// Graceful Shutdown Signal Handling //
///////////////////////////////////////

// Capture SIGINT or SIGTERM
func startSignalHandler() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go signalHandler(sigs)
}

func signalHandler(sigs chan os.Signal) {
	sig := <-sigs
	fmt.Printf("\nReceived SIG: %v\n", sig)
	gracefulShutdown()
}

func gracefulShutdown() {
	log.Println("waiting for server to gracefully shut down...")
	grpcServer.GracefulStop()
	log.Println("...server has shut down")
	os.Remove(pidFile)
	log.Println("removed lock file")
	os.Exit(0)
}
