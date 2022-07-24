package calculator

import (
	"testing"

	"github.com/rem1niscence/realtime-VWAP/subscription"
	"github.com/stretchr/testify/require"
)

func TestInvalidLimit(t *testing.T) {
	c := require.New(t)

	pair, err := NewPairVWAP(-1)
	c.Nil(pair)
	c.Error(err)
	c.Equal(err, ErrInvalidLimit)
}

func TestGetVWAP(t *testing.T) {
	c := require.New(t)

	pair, err := NewPairVWAP(1)
	c.NotNil(pair)
	c.NoError(err)

	pair.AggregateTrade("BTC-USD", &Trade{
		Price:    1,
		Quantity: 1,
	})

	c.Equal(pair.GetVWAP(), float32(1))
}

func TestAddMoreThanAllowedWindowOne(t *testing.T) {
	c := require.New(t)

	pair, err := NewPairVWAP(1)
	c.NoError(err)

	pair.AggregateTrade("BTC-USD", &Trade{
		Price:    1,
		Quantity: 1,
	})

	c.Equal(pair.GetVWAP(), float32(1))

	pair.AggregateTrade("BTC-USD", &Trade{
		Price:    2,
		Quantity: 1,
	})
	c.Equal(pair.GetVWAP(), float32(2))
}

func TestAddMoreThanAllowedWindowMultiples(t *testing.T) {
	c := require.New(t)

	pair, err := NewPairVWAP(2)
	c.NoError(err)

	pair.AggregateTrade("BTC-USD", &Trade{
		Price:    1,
		Quantity: 1,
	})
	c.Equal(pair.GetVWAP(), float32(1))

	pair.AggregateTrade("BTC-USD", &Trade{
		Price:    2,
		Quantity: 1,
	})
	c.Equal(pair.GetVWAP(), float32(1.5))

	pair.AggregateTrade("BTC-USD", &Trade{
		Price:    3,
		Quantity: 1,
	})
	c.Equal(pair.GetVWAP(), float32(2.5))
}

func TestInvalidLimitForStreamPairs(t *testing.T) {
	c := require.New(t)
	matches := make(chan subscription.Match)

	vwaps, err := StreamPairsVWAP(matches, -1)
	c.Nil(vwaps)
	c.Error(err)
}

func TestStreamPairVWAP(t *testing.T) {
	c := require.New(t)
	matches := make(chan subscription.Match)

	vwaps, err := StreamPairsVWAP(matches, 1)
	c.NotNil(vwaps)
	c.NoError(err)

	pair := "BTC-USD"
	matches <- subscription.Match{
		ProductID: pair,
		Price:     1,
		Size:      1,
	}

	vwap := <-vwaps
	c.Equal(pair, vwap.Pair)
	c.Equal(vwap.VWAP, float32(1))
}
