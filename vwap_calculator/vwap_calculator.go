package calculator

import (
	"errors"

	"github.com/rem1niscence/realtime-VWAP/subscription"
)

type Calculator interface {
	StreamPairsVWAP(matches <-chan subscription.Match, dataPoints int) <-chan *VWAP
}

var (
	// ErrInvalidLimit when a limit under 1 is given
	ErrInvalidLimit = errors.New("invalid limit value. Must be at least 1 or more")
)

type Trade struct {
	Quantity float32
	Price    float32
}

type PairVWAP struct {
	Trades                     []*Trade
	Limit                      int
	totalQuantity              float32
	totalWeightedQuantityPrice float32
	vWAP                       float32
}

type VWAP struct {
	Pair string
	VWAP float32
}

// calculateVWAP calculates, saves and retrieves the vwap value of a pair
func (pr *PairVWAP) calculateVWAP() float32 {
	pr.vWAP = pr.totalWeightedQuantityPrice / pr.totalQuantity
	return pr.vWAP
}

// Returns VWAP for a trading pair, -1 if no pair trade has been added
func (pr *PairVWAP) GetVWAP() float32 {
	return pr.vWAP
}

func (pr *PairVWAP) AggregateTrade(pair string, trade *Trade) float32 {
	if len(pr.Trades) >= pr.Limit {
		oldestPairQuantity := pr.Trades[0].Quantity
		oldestPairPrice := pr.Trades[0].Price

		if pr.Limit == 1 {
			pr.Trades = nil
		} else {
			pr.Trades = pr.Trades[1:]
		}

		pr.totalWeightedQuantityPrice -= oldestPairPrice * oldestPairQuantity
		pr.totalQuantity -= oldestPairQuantity
	}

	pr.Trades = append(pr.Trades, trade)
	pr.totalWeightedQuantityPrice += trade.Price * trade.Quantity
	pr.totalQuantity += trade.Quantity

	return pr.calculateVWAP()
}

func StreamPairsVWAP(matches <-chan subscription.Match, dataPoints int) (<-chan *VWAP, error) {
	vwapStream := make(chan *VWAP)
	tradingPairs := map[string]*PairVWAP{}

	if dataPoints <= 0 {
		return nil, ErrInvalidLimit
	}

	go func() {
		for match := range matches {
			tradingPair, ok := tradingPairs[match.ProductID]
			if !ok {
				tradingPairs[match.ProductID] = &PairVWAP{
					// As we already know beforehand the maximum size of this slice, Preallocating it will
					// add a performance boost but also create the need of having to manage the actual used
					// space for each pair in another field. IMO for this case such optimization is not needed.
					Trades: []*Trade{},
					Limit:  dataPoints,
				}
				tradingPair = tradingPairs[match.ProductID]
			}
			vwapStream <- &VWAP{
				Pair: match.ProductID,
				VWAP: tradingPair.AggregateTrade(match.ProductID, &Trade{
					Quantity: match.Size,
					Price:    match.Price,
				}),
			}
		}
	}()

	return vwapStream, nil
}
