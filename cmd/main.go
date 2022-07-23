package main

import (
	"fmt"
	"log"

	"github.com/rem1niscence/realtime-VWAP/calculator"
	"github.com/rem1niscence/realtime-VWAP/subscription"
)

func main() {
	pairs := []string{"BTC-USD", "ETH-USD", "ETH-BTC"}
	matches, err := subscription.SubscribeToMatches("wss://ws-feed.exchange.coinbase.com", pairs)
	if err != nil {
		log.Fatalf("error de todo: %s", err.Error())
	}

	vwaps := calculator.StreamPairsVWAP(matches, 200)

	for match := range vwaps {
		fmt.Printf("pair: %s | vwap: %.3f\n", match.Pair, match.VWAP)
	}
}
