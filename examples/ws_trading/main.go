// Example: WebSocket trading — connect, login, and place orders with low latency.
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
		config.WithWSReconnect(true),
	)

	wsClient.OnErrorHandler(func(err error) {
		log.Printf("WS error: %v", err)
	})

	if err := wsClient.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer wsClient.Close()

	if err := wsClient.Login(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged in successfully")

	// Enable cancel-on-disconnect for safety.
	if err := wsClient.SetCancelOnDisconnect(ctx, true); err != nil {
		log.Fatal(err)
	}

	// Place a limit order.
	order, err := wsClient.Insert(ctx,
		types.NewBuyOrderParams("BTC-PERPETUAL", 0.01).
			WithPrice(30000).
			WithPostOnly(true).
			WithLabel("ws-example"),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Order placed: %s status=%s\n", order.OrderID, order.Status)

	// Amend the order.
	amended, err := wsClient.Amend(ctx, types.NewAmendByOrderID(order.OrderID, 30100, 0.01))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Order amended: %s price=%.2f\n", amended.OrderID, *amended.Price)

	// Cancel the order.
	cancelled, err := wsClient.Cancel(ctx, types.CancelByOrderID(order.OrderID))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Order cancelled: %s\n", cancelled.OrderID)
}
