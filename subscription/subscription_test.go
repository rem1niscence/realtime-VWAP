package subscription

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

const (
	matchMsg = `{
		"type":"match",
		"trade_id":1234567689,
		"maker_order_id":"1234-1234-abcd-abcd",
		"taker_order_id":"abcd-1234-1234-abcd",
		"side":"sell",
		"size":"0.1",
		"price":"100",
		"product_id":"BTC-USD",
		"sequence":123456787,
		"time":"2022-07-24T15:03:49.910189Z"}`
)

var upgrader = websocket.Upgrader{}

func TestInvalidPairFormat(t *testing.T) {
	c := require.New(t)

	invalidPair := []string{"invalid"}
	_, err := SubscribeToMatches("", invalidPair)

	c.Error(err)
	c.Equal(err, ErrInvalidPair)
}

func TestInvalidWSConnection(t *testing.T) {
	c := require.New(t)

	pair := []string{"BTC-USD"}
	_, err := SubscribeToMatches("invalid_url", pair)
	c.Error(err)
}

func TestReadsAndForwardsMatches(t *testing.T) {
	c := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(forwarderServer([]byte(matchMsg))))
	defer server.Close()
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	pair := []string{"BTC-USD"}
	matches, err := SubscribeToMatches(wsURL, pair)
	c.NoError(err)

	match := <-matches
	c.Equal(pair[0], match.ProductID)
}

// forwarderServer will always return the same response to any websocket query.
// A better testing implementation would be to create a broadcasting server so
// custom messages can be written from the testing function but for all the purposes
// of what's being tested this will suffice.
func forwarderServer(forward []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		defer c.Close()
		for {
			mt, _, err := c.ReadMessage()
			if err != nil {
				break
			}

			err = c.WriteMessage(mt, forward)
			if err != nil {
				break
			}
		}
	}
}
