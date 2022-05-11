package websocket

import (
	"encoding/json"
	"io"
	"os"

	"net/http"

	gwebsocket "github.com/gorilla/websocket"
	"github.com/vinimmelo/zerohash/internal/logger"
	"github.com/vinimmelo/zerohash/internal/service"
	"github.com/vinimmelo/zerohash/internal/transport"
)

type WebsocketSubscriber struct {
	conn    *gwebsocket.Conn
	log     logger.Logger
	service *service.Service
}

func Dial(url string) (*gwebsocket.Conn, *http.Response, error) {
	return gwebsocket.DefaultDialer.Dial(url, nil)
}

func New(conn *gwebsocket.Conn, log logger.Logger, svc *service.Service) *WebsocketSubscriber {
	return &WebsocketSubscriber{
		conn:    conn,
		log:     log,
		service: svc,
	}
}

func (w *WebsocketSubscriber) Subscribe(request transport.SubscribeRequest) {
	data, err := json.Marshal(&request)
	if err != nil {
		w.log.Fatal("error marshalling JSON: %v", err)
	}

	err = w.conn.WriteMessage(gwebsocket.TextMessage, data)
	if err != nil {
		w.log.Fatal("error subscribing: %v", err)
	}

	for {
		mt, message, err := w.conn.ReadMessage()
		if mt == gwebsocket.CloseMessage || mt == gwebsocket.CloseNormalClosure {
			w.log.Error("server closed from the upstream")
			break
		} else if err != nil {
			w.log.Error("error while reading the message from the websocket: %v, message: %s", err, string(message))
			break
		}

		var response SubscribeResponse
		json.Unmarshal(message, &response)

		if response.Type == "subscriptions" {
			w.log.Info("subscribed successfully: %s", string(message))
			continue
		}

		tickerInfo, err := response.ToTickerInfo()
		if err != nil {
			w.log.Error("error when converting from SubscribeReponse to TickerInfo: %v, resonse: %+v, message: %s", err, response, string(message))
			continue
		}

		// Only prints to STDOUT for now, but easily extensible
		w.service.ExecuteVWAP(tickerInfo, []io.Writer{os.Stdout})
	}
}
