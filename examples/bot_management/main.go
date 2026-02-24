// Example: Bot management — creating and managing trading bots.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-thalex/auth"
	"github.com/amiwrpremium/go-thalex/config"
	"github.com/amiwrpremium/go-thalex/enums"
	"github.com/amiwrpremium/go-thalex/rest"
	"github.com/amiwrpremium/go-thalex/types"
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

	client := rest.NewClient(
		config.WithNetwork(config.Testnet),
		config.WithCredentials(creds),
	)

	endTime := float64(time.Now().Add(1 * time.Hour).Unix())

	// Create an SGSL (Signal, Go, Stop-Loss) bot.
	sgsl, err := client.CreateSGSLBot(ctx,
		types.NewSGSLBotParams("BTC-PERPETUAL", enums.TargetMark,
			45000, 0.1, 40000, 0, endTime).
			WithMaxSlippage(200).
			WithLabel("sgsl-example"),
	)
	if err != nil {
		log.Printf("SGSL bot creation: %v", err)
	} else {
		fmt.Printf("SGSL bot created: %s status=%s\n", sgsl.BotID, sgsl.Status)
	}

	// Create a Grid bot.
	grid, err := client.CreateGridBot(ctx,
		types.NewGridBotParams("BTC-PERPETUAL",
			[]float64{44000, 44500, 45000, 45500, 46000},
			0.01, endTime).
			WithLabel("grid-example"),
	)
	if err != nil {
		log.Printf("Grid bot creation: %v", err)
	} else {
		fmt.Printf("Grid bot created: %s status=%s\n", grid.BotID, grid.Status)
	}

	// Create a Delta Hedger bot.
	dhedge, err := client.CreateDHedgeBot(ctx,
		types.NewDHedgeBotParams("BTC-PERPETUAL", 60).
			WithTargetDelta(0).
			WithThreshold(0.5).
			WithEndTime(endTime).
			WithLabel("dhedge-example"),
	)
	if err != nil {
		log.Printf("DHedge bot creation: %v", err)
	} else {
		fmt.Printf("DHedge bot created: %s status=%s\n", dhedge.BotID, dhedge.Status)
	}

	// List all bots.
	bots, err := client.Bots(ctx, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nAll bots (%d):\n", len(bots))
	for _, b := range bots {
		fmt.Printf("  %s: strategy=%s status=%s instrument=%s\n",
			b.BotID, b.Strategy, b.Status, b.InstrumentName)
	}

	// Cancel all bots.
	n, err := client.CancelAllBots(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nCancelled %d bots\n", n)
}
