// Example: Basic REST API usage — fetching market data and placing orders.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-thalex/auth"
	"github.com/amiwrpremium/go-thalex/config"
	"github.com/amiwrpremium/go-thalex/rest"
	"github.com/amiwrpremium/go-thalex/types"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create an unauthenticated client for public endpoints.
	pub := rest.NewClient(
		config.WithNetwork(config.Testnet),
	)

	// Fetch active instruments.
	instruments, err := pub.Instruments(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d active instruments\n", len(instruments))

	// Fetch ticker for BTC perpetual.
	ticker, err := pub.Ticker(ctx, "BTC-PERPETUAL")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("BTC-PERPETUAL mark price: %.2f\n", ticker.MarkPrice)
	if mid := ticker.MidPrice(); mid != nil {
		fmt.Printf("BTC-PERPETUAL mid price: %.2f\n", *mid)
	}

	// Fetch order book.
	book, err := pub.Book(ctx, "BTC-PERPETUAL")
	if err != nil {
		log.Fatal(err)
	}
	if len(book.Bids) > 0 {
		fmt.Printf("Best bid: %.2f x %.4f\n", book.Bids[0].Price(), book.Bids[0].Amount())
	}
	if len(book.Asks) > 0 {
		fmt.Printf("Best ask: %.2f x %.4f\n", book.Asks[0].Price(), book.Asks[0].Amount())
	}

	// Authenticated example — requires PEM key.
	pemPath := os.Getenv("THALEX_PEM_PATH")
	keyID := os.Getenv("THALEX_KEY_ID")
	if pemPath == "" || keyID == "" {
		fmt.Println("\nSet THALEX_PEM_PATH and THALEX_KEY_ID to run authenticated examples")
		return
	}

	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		log.Fatal(err)
	}
	creds, err := auth.NewCredentialsFromPEM(keyID, pemData)
	if err != nil {
		log.Fatal(err)
	}

	client := rest.NewClient(
		config.WithNetwork(config.Testnet),
		config.WithCredentials(creds),
	)

	// Get account summary.
	summary, err := client.AccountSummary(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nAccount margin: %.2f, remaining: %.2f\n", summary.Margin, summary.RemainingMargin)

	// Place a limit buy order.
	order, err := client.Insert(ctx,
		types.NewBuyOrderParams("BTC-PERPETUAL", 0.01).
			WithPrice(30000). // Far below market to avoid fill.
			WithPostOnly(true).
			WithLabel("example-order"),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Order placed: %s status=%s\n", order.OrderID, order.Status)

	// Cancel it.
	cancelled, err := client.Cancel(ctx, types.CancelByOrderID(order.OrderID))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Order cancelled: %s status=%s\n", cancelled.OrderID, cancelled.Status)
}
