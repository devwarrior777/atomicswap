// Copyright (c) 2018/2019 The DevCo developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package xzc

import (
	"crypto/sha256"
	"net"

	"github.com/zcoinofficial/xzcd/chaincfg"
	"github.com/zcoinofficial/xzcutil"
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

// Get the default wallet port
func getWalletPort(testnet bool) string {
	if testnet {
		return "18888"
	}
	return "8888"
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

func calcFeePerKb(absoluteFee xzcutil.Amount, serializeSize int) float64 {
	return float64(absoluteFee) / float64(serializeSize) / 1e5
}
