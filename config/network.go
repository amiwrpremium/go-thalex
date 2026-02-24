// Package config provides network configuration and client options for the Thalex SDK.
//
// Use [Production] or [Testnet] to select the target environment, and
// the various With* option functions to configure client behavior.
package config

// Version is the SDK version.
const Version = "0.2.0"

// UserAgent is the default user agent sent with requests.
const UserAgent = "go-thalex/" + Version

// Network represents the Thalex network environment.
type Network int

const (
	// Production is the live trading environment.
	Production Network = iota
	// Testnet is the test trading environment.
	Testnet
)

// BaseURL returns the REST API base URL for this network.
func (n Network) BaseURL() string {
	switch n {
	case Testnet:
		return "https://testnet.thalex.com/api/v2"
	default:
		return "https://thalex.com/api/v2"
	}
}

// WebSocketURL returns the WebSocket API URL for this network.
func (n Network) WebSocketURL() string {
	switch n {
	case Testnet:
		return "wss://testnet.thalex.com/ws/api/v2"
	default:
		return "wss://thalex.com/ws/api/v2"
	}
}

// String returns a human-readable name for the network.
func (n Network) String() string {
	switch n {
	case Testnet:
		return "testnet"
	default:
		return "production"
	}
}
