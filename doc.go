// Package thalex is a Go SDK for the Thalex cryptocurrency derivatives exchange
// (https://thalex.com).
//
// The SDK is organized into the following sub-packages:
//
//   - [github.com/amiwrpremium/go-thalex/apierr] — error types (APIError, AuthError, etc.)
//   - [github.com/amiwrpremium/go-thalex/auth] — API key credentials and JWT token generation
//   - [github.com/amiwrpremium/go-thalex/config] — network selection and client configuration options
//   - [github.com/amiwrpremium/go-thalex/enums] — enum types (Direction, OrderType, etc.)
//   - [github.com/amiwrpremium/go-thalex/types] — request/response types
//   - [github.com/amiwrpremium/go-thalex/rest] — REST API client
//   - [github.com/amiwrpremium/go-thalex/ws] — WebSocket JSON-RPC client with subscriptions
//
// # Quick Start
//
//	import (
//	    "github.com/amiwrpremium/go-thalex/auth"
//	    "github.com/amiwrpremium/go-thalex/config"
//	    "github.com/amiwrpremium/go-thalex/rest"
//	)
//
//	creds, err := auth.NewCredentialsFromPEM("my-key-id", pemData)
//	client := rest.NewClient(
//	    config.WithNetwork(config.Testnet),
//	    config.WithCredentials(creds),
//	)
//	instruments, err := client.Instruments(ctx)
package thalex
