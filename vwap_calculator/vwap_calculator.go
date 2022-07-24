package calculator

import (
	"errors"

	"github.com/rem1niscence/realtime-VWAP/subscription"
)

// Calculator holds all the operations related to streaming Pair's VWAP
type Calculator interface {
	StreamPairsVWAP(matches <-chan subscription.Match, dataPoints int) <-chan *VWAP
}

var (
	// ErrInvalidLimit when a limit under 1 is given.
	ErrInvalidLimit = errors.New("invalid limit value. Must be at least 1 or more")
)

// Trade represents a trading operation.
type Trade struct {
	Quantity float32
	Price    float32
}

// PairVWAP holds the trades needed to calculate the VWAP of a pair.
type PairVWAP struct {
	trades                     []*Trade
	limit                      int
	totalQuantity              float32
	totalWeightedQuantityPrice float32
	vWAP                       float32
}

// VWAP is a helper struct to represent a pair-vwap mapping.
type VWAP struct {
	Pair string
	VWAP float32
}

// NewPairVWAP returns a PairVWAP asserting it has valid values.
func NewPairVWAP(limit int) (*PairVWAP, error) {
	if limit <= 0 {
		return nil, ErrInvalidLimit
	}

	return &PairVWAP{
		limit: limit,
	}, nil
}

// calculateVWAP calculates, saves and retrieves the vwap value of a pair.
func (pr *PairVWAP) calculateVWAP() float32 {
	pr.vWAP = pr.totalWeightedQuantityPrice / pr.totalQuantity
	return pr.vWAP
}

// GetVWAP Returns VWAP for a trading pair
func (pr *PairVWAP) GetVWAP() float32 {
	return pr.vWAP
}

// AggregateTrade adds a Trade to a pair and calculates the VWAP, dropping oldest
// values when it exceedes the window limit.
func (pr *PairVWAP) AggregateTrade(pair string, trade *Trade) float32 {
	if len(pr.trades) >= pr.limit {
		oldestPairQuantity := pr.trades[0].Quantity
		oldestPairPrice := pr.trades[0].Price

		if pr.limit == 1 {
			pr.trades = nil
		} else {
			pr.trades = pr.trades[1:]
		}

		pr.totalWeightedQuantityPrice -= oldestPairPrice * oldestPairQuantity
		pr.totalQuantity -= oldestPairQuantity
	}

	pr.trades = append(pr.trades, trade)
	pr.totalWeightedQuantityPrice += trade.Price * trade.Quantity
	pr.totalQuantity += trade.Quantity

	return pr.calculateVWAP()
}

// StreamPairsVWAP calculates VWAP for pairs given through the matches channel and
// subsequently streams them to a listener.
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
					trades: []*Trade{},
					limit:  dataPoints,
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
