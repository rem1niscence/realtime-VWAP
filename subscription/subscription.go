package subscription

import (
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
	logger "github.com/rem1niscence/realtime-VWAP/shared/log"
	"github.com/sirupsen/logrus"
)

var (
	// ErrInvalidPair when a pair with an invalid format is submitted
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

// Subscriber holds the methods relating to setting up a connection to a match channel
type Subscriber interface {
	SubscribeToMatches(address string, pairs []string) <-chan Match
}

func SubscribeToMatches(address string, pairs []string) (<-chan Match, error) {
	if arePairsValid := validatePairs(pairs); !arePairsValid {
		return nil, ErrInvalidPair
	}

	logger.Log.Debug("Connecting to:", address)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

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
		for {
			select {
			case <-interrupt:
				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				// Inspired from: https://github.com/gorilla/websocket/blob/master/examples/echo/client.go#L71
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					logger.Log.Error("ws close:", err)
					return
				}
				conn.Close()
			default:
				_, message, err := conn.ReadMessage()
				if err != nil {
					logger.Log.WithFields(logrus.Fields{
						"error": err.Error(),
					}).Error("error reading websocket message")
					return
				}

				var match Match
				err = json.Unmarshal(message, &match)
				if err != nil {
					logger.Log.WithFields(logrus.Fields{
						"error": err.Error(),
					}).Error("error parsing websocket message")
					continue
				}

				matches <- match
			}
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
