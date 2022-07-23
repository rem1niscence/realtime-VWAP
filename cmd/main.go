package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"text/tabwriter"
	"time"

	logger "github.com/rem1niscence/realtime-VWAP/pkg/log"
	"github.com/rem1niscence/realtime-VWAP/subscription"
	calculator "github.com/rem1niscence/realtime-VWAP/vwap_calculator"
	"github.com/sirupsen/logrus"
)

func main() {
	options := GatherVWAPCalculatorOptions()

	matches, err := subscription.SubscribeToMatches(options.OriginURL, options.Pairs)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("error subscribing to matches")
	}
	vwaps, err := calculator.StreamPairsVWAP(matches, options.Limit)
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
		w.Init(os.Stdout, 8, 8, 1, '\t', 0)
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
