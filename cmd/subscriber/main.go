package main

import (
	"net/url"

	"os"

	"net/http"

	"github.com/vinimmelo/zerohash/internal/logger"
	"github.com/vinimmelo/zerohash/internal/service"
	"github.com/vinimmelo/zerohash/internal/transport"
	"github.com/vinimmelo/zerohash/internal/transport/websocket"
)

func main() {
	log := logger.New(logger.INFO)
	service := service.New(log)

	host := os.Getenv("WEBSOCKET_URI")

	u := url.URL{
		Scheme: "wss",
		Host:   host,
	}

	c, _, err := websocket.Dial(u.String())
	if err != nil {
		log.Fatal("error opening a connection to: %s, err: %v", u.String(), err)
	}
	defer c.Close()
	log.Info("connected to coinbase websocket")

	http.HandleFunc("/health", HealthCheckHandler)

	go func() {
		log.Fatal("%v", http.ListenAndServe(":8080", nil))
	}()

	websocketSubscriber := websocket.New(c, log, service)
	subscribe(websocketSubscriber)
}

func subscribe(subscriber transport.Subscriber) {
	request := transport.SubscribeRequest{Type: "subscribe", ProductIds: []string{"BTC-USD", "ETH-USD", "ETH-BTC"}, Channels: []string{"ticker"}}

	subscriber.Subscribe(request)
}

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
