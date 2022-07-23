package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"text/tabwriter"
	"time"

	"github.com/rem1niscence/realtime-VWAP/calculator"
	logger "github.com/rem1niscence/realtime-VWAP/shared/log"
	"github.com/rem1niscence/realtime-VWAP/subscription"
	"github.com/sirupsen/logrus"
)

func main() {
	pairs := []string{"BTC-USD", "ETH-USD", "ETH-BTC"}
	matches, err := subscription.SubscribeToMatches("wss://ws-feed.exchange.coinbase.com", pairs)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("error subscribing to matches")
	}
	vwaps, err := calculator.StreamPairsVWAP(matches, 200)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("error setting VWAP stream")
	}

	clear, canClearScreen := clearScreen()
	if !canClearScreen {
		logger.Log.Info("Sorry, your OS does not support terminal clearing. Table results will be appended.")
	}

	pairsVWAP := map[string]float32{}
	for vwap := range vwaps {
		if canClearScreen {
			clear()
		}
		pairsVWAP[vwap.Pair] = vwap.VWAP

		w := new(tabwriter.Writer)

		// minwidth, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 1, '\t', 0)

		fmt.Println()
		fmt.Fprintf(w, "%s\t%s\t\n", "Pair", "VWAP")
		fmt.Fprintf(w, "%s\t%s\t\n", "-------", "----")

		for pair, value := range pairsVWAP {
			fmt.Fprintf(w, "%s\t%.6f\t\n", pair, value)
		}
		w.Flush()

		fmt.Println("--------------\nlast updated:", time.Now().UTC())
	}
}

func clearScreen() (func(), bool) {
	clear := make(map[string]func())

	clearUnix := func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["linux"] = clearUnix
	clear["darwin"] = clearUnix

	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	fn, ok := clear[runtime.GOOS]
	return fn, ok
}
