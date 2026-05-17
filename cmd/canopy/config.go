package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the top-level configuration for the Canopy node.
type Config struct {
	// Network settings
	ChainID    string `json:"chain_id"`
	NetworkID  uint64 `json:"network_id"`
	ListenAddr string `json:"listen_addr"`

	// RPC settings
	RPCAddr    string `json:"rpc_addr"`
	RPCEnabled bool   `json:"rpc_enabled"`

	// Consensus settings
	MaxBlockSize int64  `json:"max_block_size"`
	BlockTime    int64  `json:"block_time_ms"`

	// Storage settings
	DataDir string `json:"data_dir"`
	LogLevel string `json:"log_level"`

	// Validator settings
	ValidatorKey string `json:"validator_key,omitempty"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return &Config{
		ChainID:      "canopy-1",
		NetworkID:    1,
		ListenAddr:   "0.0.0.0:9090",
		RPCAddr:      "0.0.0.0:50832",
		RPCEnabled:   true,
		MaxBlockSize: 4 * 1024 * 1024, // 4 MB
		BlockTime:    4000,             // 4 seconds
		DataDir:      filepath.Join(homeDir, ".canopy"),
		LogLevel:     "info",
	}
}

// LoadConfig reads a JSON config file from the given path and merges it
// on top of the default configuration.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return defaults when no config file exists yet.
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	return cfg, nil
}

// Save writes the configuration to disk as pretty-printed JSON,
// creating any necessary parent directories.
func (c *Config) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling config: %w", err)
	}

	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("writing config file %q: %w", path, err)
	}

	return nil
}

// Validate performs basic sanity checks on the configuration.
func (c *Config) Validate() error {
	if c.ChainID == "" {
		return fmt.Errorf("chain_id must not be empty")
	}
	if c.DataDir == "" {
		return fmt.Errorf("data_dir must not be empty")
	}
	if c.BlockTime <= 0 {
		return fmt.Errorf("block_time_ms must be positive")
	}
	if c.MaxBlockSize <= 0 {
		return fmt.Errorf("max_block_size must be positive")
	}
	return nil
}
