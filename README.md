# Realtime VWAP

Calculator of [VWAP (Volume-weighted average price)](https://en.wikipedia.org/wiki/Volume-weighted_average_price) values for specified trading pairs using realtime values from Coinbase

## Design

The design of this is quite simple but functional, work is splitted mainly between two goroutines which run in a top-down fashion:
The first goroutine is inside the `SubscribeToMatches` function, which connects to the Coinbase API and it keeps listening to any
upcoming trade for the specified pairs, exposing a channel on which other goroutines can listen to the results obtained.
From there on the second main goroutine inside `StreamPairsVWAP` takes the channel exposed by the previous function and proceeds to
calculate and save the VWAP of the given trade. This function also exposes a channel to listen to the VWAP results once they're calculated and this last channel is used by the `main` goroutine to print the results to `std`. This design does not make any
assumption more than the default environment values in the [options.go](cmd/options.go).

## How to run this

There are up to 4 different ways to run this program:

### Option 1: Executable

There's a executable of the latest version in the project's root, assuming you're on a *NIX terminal you need to run the following command:

```console
$ ./realtime_vwap

Pair    VWAP
------- ----
BTC-USD 22710.613281
ETH-USD 1597.085205
ETH-BTC 0.070340
--------------
Last updated: 2022-07-24 18:34:12.331341 +0000 UTC
```

Bear in mind this executable was built for *NIX systems so running this on windows or any other type of OS may not work.

### Option 2: Building it yourself

You can build it yourself if you have a go compiler in your computer for that you need to run from the project's root folder `go build -o realtime_vwap cmd/*.go` and run the executable just like the previous step. If you don't want to run it directly you can also run `go run cmd/*.go` whih won't generate any executable.

### Option 3: docker-compose

Assuming you have [docker-compose](https://docs.docker.com/compose/install/) intalled, you only need to run `docker-compose up` in your terminal and the program will run.

### Option 4: Docker

Assuming you have [docker](https://docs.docker.com/get-docker/) installed, you need to run the following commands in your terminal:

```console
$ docker build -t realtime_vwap .
Build output...

$ docker run realtime_vwap
```

And that's it, the program will be shown in your terminal.
