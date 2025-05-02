package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ekival-labs/ekival-canvas/config"
	"github.com/ekival-labs/ekival-canvas/router"
)

var cmdlineFlags struct {
	configFile string
}

// var httpClient *http.Client // Global HTTP client with proxy setup

// func init() {
// 	// Configure the SOCKS5 proxy
// 	proxyAddr := "127.0.0.1:2080"
// 	proxyDialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
// 	if err != nil {
// 		log.Fatalf("Failed to create SOCKS5 dialer: %v", err)
// 	}

// 	// Set up the HTTP transport with the proxy
// 	transport := &http.Transport{
// 		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
// 			return proxyDialer.Dial(network, addr)
// 		},
// 	}

// 	// Set the default HTTP transport globally
// 	http.DefaultTransport = transport

// 	log.Println("Proxy initialized globally.")
// }

func main() {

	// Load config
	flag.StringVar(
		&cmdlineFlags.configFile,
		"config",
		"",
		"path to config file to load",
	)
	flag.Parse()

	err := config.Load(cmdlineFlags.configFile)
	if err != nil {
		fmt.Printf("Failed to load config: %s\n", err)
		os.Exit(1)
	}

	err = config.ChainCTXSetup()
	if err != nil {
		fmt.Printf("Failed to set up chain context: %s\n", err)
		os.Exit(1)
	}

	// Set up other configurations
	config.WalletSetup()

	// Initialize router
	router.RouterInit()
}
