// Package main is the entry point for the Canopy node application.
// Canopy is a proof-of-stake blockchain network implementation in Go.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/canopy-network/canopy/lib"
	"github.com/canopy-network/canopy/lib/crypto"
	"github.com/canopy-network/canopy/node"
)

const (
	// AppName is the name of the application
	AppName = "canopy"
	// AppVersion is the current version of the application
	AppVersion = "0.0.1"
)

func main() {
	// Parse CLI arguments and configuration
	config, err := lib.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize the logger
	logger := lib.NewLogger(config.LogLevel)
	logger.Infof("Starting %s v%s", AppName, AppVersion)

	// Load or generate node private key
	privKey, err := crypto.LoadOrGeneratePrivateKey(config.DataDir)
	if err != nil {
		logger.Errorf("failed to load or generate private key: %v", err)
		os.Exit(1)
	}
	logger.Infof("Node public key: %s", privKey.PublicKey().String())

	// Initialize and start the node
	n, err := node.NewNode(config, privKey, logger)
	if err != nil {
		logger.Errorf("failed to create node: %v", err)
		os.Exit(1)
	}

	if err := n.Start(); err != nil {
		logger.Errorf("failed to start node: %v", err)
		os.Exit(1)
	}

	logger.Info("Node started successfully")

	// Wait for interrupt signal to gracefully shut down the node
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down node...")
	if err := n.Stop(); err != nil {
		logger.Errorf("error during node shutdown: %v", err)
		os.Exit(1)
	}

	logger.Info("Node stopped gracefully")
}
