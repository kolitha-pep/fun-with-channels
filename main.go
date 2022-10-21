package main

import (
	"github.com/kolitha-pep/fun-with-channels/app/handler/multiCurrencySma"
	"github.com/kolitha-pep/fun-with-channels/app/pkg/finnhub"
)

const FinnWebsocketURL = "wss://ws.finnhub.io?token=ccpcv5qad3i91ts8jhk0ccpcv5qad3i91ts8jhkg"

func main() {

	// connect to finnhub websocket
	ws, err := finnhub.WebsocketDialer(FinnWebsocketURL)
	if err != nil {
		panic(err)
	}
	defer ws.Close()

	//singleCurrencySma.NewSimpleMovingAverage(ws).Calculate()
	multiCurrencySma.NewSimpleMovingAverage(ws).Calculate()
}
