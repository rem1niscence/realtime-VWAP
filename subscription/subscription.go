package subscription

import (
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

// Subscriber holds the methods relating to setting up a connection to a match channel.
type Subscriber interface {
	SubscribeToMatches(address string, pairs []string) <-chan Match
}

var (
	// ErrInvalidPair when a pair with an invalid format is submitted.
	ErrInvalidPair = errors.New("subscriber: invalid pair format. Format must be '{min 3 uppercase}-{min 3 uppercase}'")
)

// Match represents a trading pair match
type Match struct {
	Type         string    `json:"type"`
	TradeID      int       `json:"trade_id"`
	MakerOrderID string    `json:"maker_order_id"`
	TakerOrderID string    `json:"taker_order_id"`
	Side         string    `json:"side"`
	Size         float32   `json:"size,string"`
	Price        float32   `json:"price,string"`
	ProductID    string    `json:"product_id"`
	Sequence     int64     `json:"sequence"`
	Time         time.Time `json:"time"`
}

// SubscribeToMatches returns a channel to listen to trading operations for the given pairs.
func SubscribeToMatches(address string, pairs []string) (<-chan Match, error) {
	if arePairsValid := validatePairs(pairs); !arePairsValid {
		return nil, ErrInvalidPair
	}

	conn, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		return nil, err
	}

	msg := struct {
		Type       string   `json:"type"`
		Channels   []string `json:"channels"`
		ProductIDs []string `json:"product_ids"`
	}{
		Type:       "subscribe",
		Channels:   []string{"matches"},
		ProductIDs: pairs,
	}
	rawMsg, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}

	if err := conn.WriteMessage(1, rawMsg); err != nil {
		return nil, err
	}

	matches := make(chan Match)
	go func() {
		defer close(matches)
		defer conn.Close()
		for {
			var match Match
			err := conn.ReadJSON(&match)
			if err != nil {
				log.Printf("read websocket message: %s\n", err.Error())
				continue
			}

			// Invalid match, do not send.
			if len(match.ProductID) == 0 {
				continue
			}

			matches <- match
		}
	}()

	return matches, nil
}

// validatePairs asserts whether all the pairs to match come in the form `{min 3 alphanumeric}-{min 3 alphanumeric}`.
func validatePairs(pairs []string) bool {
	r := regexp.MustCompile(`[A-Z0-9]{3,}-[A-Z0-9]{3,}`)
	for _, pair := range pairs {
		if match := r.MatchString(pair); !match {
			return false
		}
	}
	return true
}
