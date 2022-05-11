package websocket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/vinimmelo/zerohash/internal/entities"
	"github.com/vinimmelo/zerohash/internal/logger"
	"github.com/vinimmelo/zerohash/internal/service"
	"github.com/vinimmelo/zerohash/internal/transport"
)

var upgrader = websocket.Upgrader{}

var msg1 string = `{
   "type":"ticker",
   "sequence":29208664515,
   "product_id":"ETH-USD",
   "price":"2299.23",
   "open_24h":"2362.2",
   "volume_24h":"399462.93411513",
   "low_24h":"2113.79",
   "high_24h":"2451.12",
   "volume_30d":"4967243.37945517",
   "best_bid":"2298.63",
   "best_ask":"2299.23",
   "side":"buy",
   "time":"2022-05-11T15:24:41.016034Z",
   "trade_id":271814604,
   "last_size":"0.0246174"
}`

var msg2 string = `{
   "type":"ticker",
   "sequence":37372821076,
   "product_id":"BTC-USD",
   "price":"31247.72",
   "open_24h":"31320.22",
   "volume_24h":"43243.32954425",
   "low_24h":"29026.26",
   "high_24h":"32185.9",
   "volume_30d":"532652.23995953",
   "best_bid":"31241.92",
   "best_ask":"31247.73",
   "side":"sell",
   "time":"2022-05-11T15:24:40.635431Z",
   "trade_id":332372454,
   "last_size":"0.02351763"
}`

var msg3 string = `{
   "type":"ticker",
   "sequence":5336071653,
   "product_id":"ETH-BTC",
   "price":"0.07355",
   "open_24h":"0.07543",
   "volume_24h":"14434.00695329",
   "low_24h":"0.07295",
   "high_24h":"0.0769",
   "volume_30d":"113231.29909807",
   "best_bid":"0.07354",
   "best_ask":"0.07355",
   "side":"buy",
   "time":"2022-05-11T15:24:39.052541Z",
   "trade_id":27360499,
   "last_size":"0.002"
} `

var msg4 string = `{
   "type":"ticker",
   "sequence":5336071653,
   "product_id":"ETH-BTC",
   "price":"0.07353",
   "open_24h":"0.07543",
   "volume_24h":"14434.00695329",
   "low_24h":"0.07295",
   "high_24h":"0.0769",
   "volume_30d":"113231.29909807",
   "best_bid":"0.07354",
   "best_ask":"0.07355",
   "side":"buy",
   "time":"2022-05-11T15:24:39.052541Z",
   "trade_id":27360499,
   "last_size":"0.002"
} `

type WebsocketServer struct {
	t *testing.T
}

func (ws WebsocketServer) mockCoinbaseResponse(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()

	mt, _, err := c.ReadMessage()
	if err != nil || mt == websocket.CloseMessage {
		return
	}

	messages := []string{`{"type": "subscriptions"}`, msg1, msg2, msg3, msg4}
	for _, msg := range messages {
		ws.t.Logf("message writed: %s", msg)
		err = c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			break
		}
	}

	c.WriteMessage(websocket.CloseNormalClosure, []byte{})
}

func Test_WebsocketSubscribe(t *testing.T) {
	ws := WebsocketServer{t: t}
	s := httptest.NewServer(http.HandlerFunc(ws.mockCoinbaseResponse))
	defer s.Close()

	url := "ws" + strings.TrimPrefix(s.URL, "http")
	t.Logf("url: %s", url)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	log := logger.New(logger.DEBUG)
	svc := service.New(log)

	websocketSubscriber := New(c, log, svc)
	request := transport.SubscribeRequest{Type: "subscribe", ProductIds: []string{"BTC-USD", "ETH-USD", "ETH-BTC"}, Channels: []string{"ticker"}}

	t.Log("Dependencies initiated")

	websocketSubscriber.Subscribe(request)

	t.Log("Subscribed successfully")

	btcUsdQueue := service.Queues[entities.BTCUSD]
	ethBtcQueue := service.Queues[entities.ETHBTC]
	ethUsdQueue := service.Queues[entities.ETHUSD]

	if len(btcUsdQueue.Elements()) != 1 {
		t.Errorf("BTC-USD queue should have 1 element")
	}

	if len(ethBtcQueue.Elements()) != 2 {
		t.Errorf("ETH-BTC queue should have 2 element")
	}

	if len(ethUsdQueue.Elements()) != 1 {
		t.Errorf("ETH-USD queue should have 1 element")
	}
}
