package singleCurrencySma

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kolitha-pep/fun-with-channels/app/pkg/datastore"
	"github.com/kolitha-pep/fun-with-channels/app/pkg/errors"
	"github.com/kolitha-pep/fun-with-channels/app/pkg/finnhub"
	"github.com/kolitha-pep/fun-with-channels/pkg/utils"
)

const (
	ResponseTypeTrade = "trade"
	WindowSize        = 10
)

type simpleMovingAverage struct {
	LastPrices          []float64
	Window              int
	SimpleMovingAverage float64
	ws                  *websocket.Conn
	sync.Mutex
}

type SimpleMovingAverageInterface interface {
	Calculate()
}

func NewSimpleMovingAverage(w *websocket.Conn) SimpleMovingAverageInterface {
	return &simpleMovingAverage{
		ws: w,
	}
}

func (t *simpleMovingAverage) Calculate() {
	symbols := []string{"BINANCE:BTCUSDT"}
	for _, s := range symbols {
		msg, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})
		t.ws.WriteMessage(websocket.TextMessage, msg)
	}

	// create a channel and pass it to process function
	sme := make(chan interface{})
	defer close(sme)
	go t.process(sme)
	var msg interface{}

	for {
		err := t.ws.ReadJSON(&msg)
		if err != nil {
			panic(err)
		}
		sme <- msg
	}
}

func (t *simpleMovingAverage) process(ch chan interface{}) {
	fmt.Println("Start calculating simple moving average")
	t.Window = WindowSize

	iteration := 0
	for {
		t.Lock()
		iteration += 1
		chanData := <-ch

		response := utils.FinnRespToStruct(chanData)
		respAverage, err := getRespAverage(response)
		if err != nil {
			log.Print(err, " skipping...")
			errors.Log(err)
			t.Unlock()
			continue
		}

		t.LastPrices = append(t.LastPrices, respAverage)

		if response.Type == ResponseTypeTrade {
			t.calcSimpleMovingAverage()
			if iteration%t.Window == 0 {
				out := fmt.Sprintf("SMA [\"BINANCE:BTCUSDT\"]: %f time: %s \n", t.SimpleMovingAverage, time.Now().String())
				err := datastore.WriteFile(out, "sma_records.txt")
				if err != nil {
					log.Print(err, " skipping...")
					errors.Log(fmt.Errorf("%s write data to file: %s \n", time.Now().String(), err.Error()))
					t.Unlock()
					continue
				}
				fmt.Println("SMA [\"BINANCE:BTCUSDT\"]: ", t.SimpleMovingAverage)
			}
		}
		t.Unlock()
	}
}

func getRespAverage(resp finnhub.FinnResp) (float64, error) {
	if len(resp.Data) == 0 {
		return 0, fmt.Errorf("%s get response average: no response data \n", time.Now().String())
	}

	var sum float64
	for _, d := range resp.Data {
		sum += d.P
	}
	return sum / float64(len(resp.Data)), nil
}

func (t *simpleMovingAverage) calcSimpleMovingAverage() {
	var sum float64
	for _, p := range t.LastPrices {
		sum += p
	}
	t.SimpleMovingAverage = sum / float64(len(t.LastPrices))
	t.LastPrices = nil
}
