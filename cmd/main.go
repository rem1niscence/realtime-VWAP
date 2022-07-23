package main

import (
	"fmt"
	"log"

	"github.com/rem1niscence/realtime-VWAP/subscription"
)

func main() {
	pairs := []string{"BTC-USD", "ETH-USD", "ETH-BTC"}
	matches, err := subscription.SubscribeToMatches("wss://ws-feed.exchange.coinbase.com", pairs)
	if err != nil {
		log.Fatalf("error de todo: %s", err.Error())
	}

	for match := range matches {
		fmt.Printf("%+v\n", match)
	}
}
