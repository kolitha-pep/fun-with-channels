package finnhub

import (
	"github.com/gorilla/websocket"
)

type FinnResp struct {
	Type string         `json:"type"`
	Data []FinnRespData `json:"data"`
}

type FinnRespData struct {
	S string  `json:"s"`
	P float64 `json:"p"`
	T int64   `json:"t"`
	V float64 `json:"v"`
}

func WebsocketDialer(url string) (*websocket.Conn, error) {
	w, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return w, err
	}
	return w, nil
}
