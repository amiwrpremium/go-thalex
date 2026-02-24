//go:build integration

package rest_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/amiwrpremium/go-thalex/auth"
	"github.com/amiwrpremium/go-thalex/config"
	"github.com/amiwrpremium/go-thalex/rest"
)

func setupIntegrationClient(t *testing.T) *rest.Client {
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

	return rest.NewClient(
		config.WithNetwork(config.Testnet),
		config.WithCredentials(creds),
	)
}

func TestIntegration_Instruments(t *testing.T) {
	client := setupIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	instruments, err := client.Instruments(ctx)
	if err != nil {
		t.Fatalf("Instruments: %v", err)
	}
	if len(instruments) == 0 {
		t.Fatal("expected at least one instrument")
	}

	// Verify BTC-PERPETUAL exists.
	found := false
	for _, inst := range instruments {
		if inst.InstrumentName == "BTC-PERPETUAL" {
			found = true
			break
		}
	}
	if !found {
		t.Error("BTC-PERPETUAL not found in instruments")
	}
}

func TestIntegration_Ticker(t *testing.T) {
	client := setupIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	ticker, err := client.Ticker(ctx, "BTC-PERPETUAL")
	if err != nil {
		t.Fatalf("Ticker: %v", err)
	}
	if ticker.MarkPrice <= 0 {
		t.Errorf("expected positive mark price, got %f", ticker.MarkPrice)
	}
}

func TestIntegration_Index(t *testing.T) {
	client := setupIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	idx, err := client.Index(ctx, "BTCUSD")
	if err != nil {
		t.Fatalf("Index: %v", err)
	}
	if idx.Price <= 0 {
		t.Errorf("expected positive index price, got %f", idx.Price)
	}
}

func TestIntegration_Book(t *testing.T) {
	client := setupIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	book, err := client.Book(ctx, "BTC-PERPETUAL")
	if err != nil {
		t.Fatalf("Book: %v", err)
	}
	if len(book.Bids) == 0 && len(book.Asks) == 0 {
		t.Log("warning: book has no bids or asks (may be normal on testnet)")
	}
}

func TestIntegration_SystemInfo(t *testing.T) {
	client := setupIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	info, err := client.SystemInfo(ctx)
	if err != nil {
		t.Fatalf("SystemInfo: %v", err)
	}
	if info.Version == "" {
		t.Error("expected non-empty system version")
	}
}

func TestIntegration_Portfolio(t *testing.T) {
	client := setupIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := client.Portfolio(ctx)
	if err != nil {
		t.Fatalf("Portfolio: %v", err)
	}
}

func TestIntegration_AccountSummary(t *testing.T) {
	client := setupIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	summary, err := client.AccountSummary(ctx)
	if err != nil {
		t.Fatalf("AccountSummary: %v", err)
	}
	// Account should exist.
	_ = summary
}

func TestIntegration_OpenOrders(t *testing.T) {
	client := setupIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	orders, err := client.OpenOrders(ctx)
	if err != nil {
		t.Fatalf("OpenOrders: %v", err)
	}
	_ = orders // May be empty.
}

func TestIntegration_Bots(t *testing.T) {
	client := setupIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	bots, err := client.Bots(ctx, true)
	if err != nil {
		t.Fatalf("Bots: %v", err)
	}
	_ = bots // May be empty.
}

func TestIntegration_ConditionalOrders(t *testing.T) {
	client := setupIntegrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	orders, err := client.ConditionalOrders(ctx)
	if err != nil {
		t.Fatalf("ConditionalOrders: %v", err)
	}
	_ = orders // May be empty.
}
