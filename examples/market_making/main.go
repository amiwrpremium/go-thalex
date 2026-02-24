// Example: Market making — mass quotes and MM protection via WebSocket.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-thalex/auth"
	"github.com/amiwrpremium/go-thalex/config"
	"github.com/amiwrpremium/go-thalex/types"
	"github.com/amiwrpremium/go-thalex/ws"
)

func main() {
	pemPath := os.Getenv("THALEX_PEM_PATH")
	keyID := os.Getenv("THALEX_KEY_ID")
	if pemPath == "" || keyID == "" {
		log.Fatal("Set THALEX_PEM_PATH and THALEX_KEY_ID environment variables")
	}

	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		log.Fatal(err)
	}
	creds, err := auth.NewCredentialsFromPEM(keyID, pemData)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	wsClient := ws.NewClient(
		config.WithNetwork(config.Testnet),
		config.WithCredentials(creds),
	)

	if err := wsClient.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer wsClient.Close()

	if err := wsClient.Login(ctx); err != nil {
		log.Fatal(err)
	}

	// Set cancel-on-disconnect for safety.
	wsClient.SetCancelOnDisconnect(ctx, true)

	// Configure market maker protection.
	err = wsClient.SetMMProtection(ctx, &types.MMProtectionParams{
		Product:     "FBTCUSD",
		TradeAmount: 5.0,
		QuoteAmount: 25.0,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MM protection configured")

	// Send a mass quote with single-level quotes.
	result, err := wsClient.MassQuote(ctx,
		types.NewMassQuoteParams([]types.DoubleSidedQuote{
			types.NewSingleLevelQuote("BTC-PERPETUAL", 30000, 0.1, 80000, 0.1),
		}).WithPostOnly(true).WithLabel("mm-example"),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Mass quote: %d success, %d fail\n", result.NSuccess, result.NFail)
	for _, e := range result.Errors {
		fmt.Printf("  Error: %s (code %d)\n", e.Message, e.Code)
	}

	// Send multi-level quotes.
	result, err = wsClient.MassQuote(ctx,
		types.NewMassQuoteParams([]types.DoubleSidedQuote{
			types.NewDoubleSidedQuote("BTC-PERPETUAL",
				[]types.QuoteLevel{
					{Price: 30000, Amount: 0.05},
					{Price: 29900, Amount: 0.10},
				},
				[]types.QuoteLevel{
					{Price: 80000, Amount: 0.05},
					{Price: 80100, Amount: 0.10},
				},
			),
		}).WithPostOnly(true),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Multi-level quote: %d success, %d fail\n", result.NSuccess, result.NFail)

	// Cancel all mass quotes.
	if err := wsClient.CancelMassQuote(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("All mass quotes cancelled")
}
