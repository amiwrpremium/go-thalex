// Example: Real-time subscriptions — streaming tickers, order books, and account updates.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/amiwrpremium/go-thalex/auth"
	"github.com/amiwrpremium/go-thalex/config"
	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/types"
	"github.com/amiwrpremium/go-thalex/ws"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	wsClient := ws.NewClient(
		config.WithNetwork(config.Testnet),
		config.WithWSReconnect(true),
	)

	wsClient.OnErrorHandler(func(err error) {
		log.Printf("WS error: %v", err)
	})

	if err := wsClient.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer wsClient.Close()

	// Subscribe to BTC ticker.
	tickerCh := types.TickerChannel("BTC-PERPETUAL", enums.Delay100ms)
	wsClient.OnTicker(tickerCh, func(t types.Ticker) {
		fmt.Printf("[ticker] BTC mark=%.2f", t.MarkPrice)
		if t.BestBidPrice != nil && t.BestAskPrice != nil {
			fmt.Printf(" bid=%.2f ask=%.2f", *t.BestBidPrice, *t.BestAskPrice)
		}
		fmt.Println()
	})

	// Subscribe to ETH lightweight ticker.
	lwtCh := types.LWTChannel("ETH-PERPETUAL", enums.Delay100ms)
	wsClient.OnLWT(lwtCh, func(t types.LightweightTicker) {
		fmt.Printf("[lwt] ETH mark=%.2f\n", t.MarkPrice)
	})

	// Subscribe to BTC order book.
	bookCh := types.BookChannel("BTC-PERPETUAL", 1, 5, enums.Delay100ms)
	wsClient.OnBook(bookCh, func(b types.BookUpdate) {
		fmt.Printf("[book] %d bids, %d asks, %d trades\n",
			len(b.Bids), len(b.Asks), len(b.Trades))
	})

	// Subscribe to index price.
	idxCh := types.PriceIndexChannel("BTCUSD")
	wsClient.OnPriceIndex(idxCh, func(idx types.IndexPrice) {
		fmt.Printf("[index] %s = %.2f\n", idx.IndexName, idx.Price)
	})

	// Subscribe to instrument changes.
	wsClient.OnInstruments(func(instruments []types.Instrument) {
		fmt.Printf("[instruments] %d instruments\n", len(instruments))
	})

	if err := wsClient.Subscribe(ctx, tickerCh, lwtCh, bookCh, idxCh, types.ChannelInstruments); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Subscribed to public channels. Press Ctrl+C to exit.")

	// Optionally subscribe to private channels if credentials are available.
	pemPath := os.Getenv("THALEX_PEM_PATH")
	keyID := os.Getenv("THALEX_KEY_ID")
	if pemPath != "" && keyID != "" {
		pemData, err := os.ReadFile(pemPath)
		if err != nil {
			log.Fatal(err)
		}
		creds, err := auth.NewCredentialsFromPEM(keyID, pemData)
		if err != nil {
			log.Fatal(err)
		}

		// Re-create with credentials for login.
		_ = creds
		loginCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// We need to create a new client with credentials to login.
		// For simplicity, re-use the existing one by generating a token.
		token, err := creds.GenerateToken()
		if err != nil {
			log.Fatal(err)
		}
		_ = token

		// In practice, create the WSClient with WithCredentials from the start,
		// then call wsClient.Login(ctx).
		fmt.Println("(Private subscription example requires creating client with credentials)")
		_ = loginCtx
	}

	<-ctx.Done()
	fmt.Println("\nShutting down...")
}
