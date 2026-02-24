//go:build integration

package ws_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/amiwrpremium/go-thalex/auth"
	"github.com/amiwrpremium/go-thalex/config"
	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
	"github.com/amiwrpremium/go-thalex/ws"
)

func setupWSIntegrationClient(t *testing.T) *ws.Client {
	t.Helper()

	pemPath := os.Getenv("THALEX_PEM_PATH")
	keyID := os.Getenv("THALEX_KEY_ID")

	if pemPath == "" || keyID == "" {
		t.Skip("THALEX_PEM_PATH and THALEX_KEY_ID must be set for integration tests")
	}

	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		t.Fatalf("reading PEM file: %v", err)
	}

	creds, err := auth.NewCredentialsFromPEM(keyID, pemData)
	if err != nil {
		t.Fatalf("creating credentials: %v", err)
	}

	return ws.NewClient(
		config.WithNetwork(config.Testnet),
		config.WithCredentials(creds),
	)
}

func TestIntegration_WS_ConnectAndLogin(t *testing.T) {
	client := setupWSIntegrationClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer client.Close()

	if !client.IsConnected() {
		t.Fatal("expected IsConnected to be true after Connect")
	}

	if err := client.Login(ctx); err != nil {
		t.Fatalf("Login: %v", err)
	}
}

func TestIntegration_WS_PublicSubscription(t *testing.T) {
	client := setupWSIntegrationClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer client.Close()

	var mu sync.Mutex
	var received bool

	ch := types.TickerChannel("BTC-PERPETUAL", enums.Delay100ms)
	client.OnTicker(ch, func(t types.Ticker) {
		mu.Lock()
		received = true
		mu.Unlock()
	})

	if err := client.Subscribe(ctx, ch); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	// Wait for at least one ticker update.
	deadline := time.After(15 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			t.Fatal("timed out waiting for ticker update")
		case <-ticker.C:
			mu.Lock()
			got := received
			mu.Unlock()
			if got {
				return // Success.
			}
		}
	}
}

func TestIntegration_WS_SetCancelOnDisconnect(t *testing.T) {
	client := setupWSIntegrationClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer client.Close()

	if err := client.Login(ctx); err != nil {
		t.Fatalf("Login: %v", err)
	}

	if err := client.SetCancelOnDisconnect(ctx, true); err != nil {
		t.Fatalf("SetCancelOnDisconnect: %v", err)
	}
}
