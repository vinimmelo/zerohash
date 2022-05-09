package main

import (
	"fmt"
	"log"

	"encoding/json"
	"net/url"

	"github.com/gorilla/websocket"
)

type SubscribeRequest struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

func main() {

	u := url.URL{
		Scheme: "wss",
		Host:   "ws-feed.exchange.coinbase.com",
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("error opening a connection to: %s, err: %v", u.String(), err)
	}
	defer c.Close()

	subscribeRequest := SubscribeRequest{Type: "subscribe", ProductIds: []string{"BTC-USD", "ETH-USD", "ETH-BTC"}, Channels: []string{"ticker"}}
	data, err := json.Marshal(&subscribeRequest)
	if err != nil {
		log.Fatalf("error marshalling JSON, %v", err)
	}

	fmt.Printf("data: %s", string(data))
	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Fatalf("error writing message: %v", err)
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
		}
		fmt.Printf("message received: %v\n", string(message))
	}
}
