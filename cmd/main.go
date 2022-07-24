package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"text/tabwriter"
	"time"

	"github.com/rem1niscence/realtime-VWAP/subscription"
	calculator "github.com/rem1niscence/realtime-VWAP/vwap_calculator"
)

func main() {
	options := GatherVWAPCalculatorOptions()

	matches, err := subscription.SubscribeToMatches(options.OriginURL, options.Pairs)
	if err != nil {
		log.Fatalf("subscribe to matches: %s\n", err.Error())
	}
	vwaps, err := calculator.StreamPairsVWAP(matches, options.Limit)
	if err != nil {
		log.Fatalf("VWAP stream: %s\n", err.Error())
	}

	clear, canClearScreen := clearScreen()
	if !canClearScreen {
		log.Println("Sorry, your OS does not support terminal clearing. Table results will be appended.")
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

		fmt.Println("--------------\nLast updated:", time.Now().UTC())
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
