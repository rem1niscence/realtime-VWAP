package main

import (
	"strings"

	"github.com/rem1niscence/realtime-VWAP/pkg/environment"
)

const (
	// OriginURL is the URL to connect in order to retrieve the data.
	OriginURL = "ORIGIN_URL"
	// DataPoints is the maximum allowed number of data points to take into account to calculate the VWAP.
	DataPoints = "DATA_POINTS"
	// Pairs is the series of pairs to look for and calculate VWAP.
	Pairs = "PAIRS"
)

// VWAPCalculatorOptions are the options needed in order to run the program.
type VWAPCalculatorOptions struct {
	OriginURL string
	Limit     int
	Pairs     []string
}

// GatherVWAPCalculatorOptions retrieves all the needed config from environment variables to a Golang struct
func GatherVWAPCalculatorOptions() VWAPCalculatorOptions {
	return VWAPCalculatorOptions{
		OriginURL: environment.GetString(OriginURL, "wss://ws-feed.exchange.coinbase.com"),
		Limit:     int(environment.GetInt64(DataPoints, 200)),
		Pairs:     strings.Split(environment.GetString(Pairs, "BTC-USD,ETH-USD,ETH-BTC"), ","),
	}
}
