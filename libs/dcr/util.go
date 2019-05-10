// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package dcr

import (
	"crypto/sha256"
	"net"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/dcrutil"
)

// Get a normalized address from `addr' which can be of form Host[:Port]`
func getNormalizedAddress(testnet bool, addr string) (hostport string, err error) {
	host, port, origErr := net.SplitHostPort(addr)
	if origErr == nil {
		return net.JoinHostPort(host, port), nil
	}
	defaultPort := getWalletPort(testnet)
	addr = net.JoinHostPort(addr, defaultPort)
	_, _, err = net.SplitHostPort(addr)
	if err != nil {
		return "", origErr
	}
	return addr, nil
}

// Get the default GPRC wallet port
func getWalletPort(testnet bool) string {
	if testnet {
		return "19111"
	}
	return "9111"
}

// Get all of the chain parameters for a network
func getChainParams(testnet bool) *chaincfg.Params {
	if testnet {
		return &chaincfg.TestNet3Params
	}
	return &chaincfg.MainNetParams
}

func sha256Hash(x []byte) []byte {
	h := sha256.Sum256(x)
	return h[:]
}

func calcFeePerKb(absoluteFee dcrutil.Amount, serializeSize int) float64 {
	return float64(absoluteFee) / float64(serializeSize) / 1e5
}

// reverse between byteslices representing wire and 'normal' hashes, ids
func byteRev(in []byte) []byte {
	inLen := len(in)
	if inLen == 0 {
		return in
	}
	out := make([]byte, inLen)
	for i := 0; i < inLen; i++ {
		out[i] = in[inLen-i-1]
	}
	return out
}
