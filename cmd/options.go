package main

import (
	"strings"

	"github.com/rem1niscence/realtime-VWAP/pkg/environment"
)

const (
	ORIGIN_URL  = "ORIGIN_URL"
	DATA_POINTS = "DATA_POINTS"
	PAIRS       = "PAIRS"
)

type VWAPCalculatorOptions struct {
	OriginURL string
	Limit     int
	Pairs     []string
}

func GatherVWAPCalculatorOptions() VWAPCalculatorOptions {
	return VWAPCalculatorOptions{
		OriginURL: environment.GetString(ORIGIN_URL, "wss://ws-feed.exchange.coinbase.com"),
		Limit:     int(environment.GetInt64(DATA_POINTS, 200)),
		Pairs:     strings.Split(environment.GetString(PAIRS, "BTC-USD,ETH-USD,ETH-BTC"), ","),
	}
}
