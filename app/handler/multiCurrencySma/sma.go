package multiCurrencySma

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

var symbols = []string{"BINANCE:BTCUSDT", "BINANCE:ETHUSDT", "BINANCE:ADAUSDT"}

type simpleMovingAverage struct {
	LastPrices          map[string][]float64
	Window              int
	SimpleMovingAverage map[string]float64
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
	for _, s := range symbols {
		msg, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})
		t.ws.WriteMessage(websocket.TextMessage, msg)
	}

	// create a channel and pass it to process function
	sme := make(chan []finnhub.FinnRespData)
	defer close(sme)
	go t.process(sme)
	var msg interface{}

	for {
		err := t.ws.ReadJSON(&msg)
		if err != nil {
			panic(err)
		}
		response := utils.FinnRespToStruct(msg)
		if response.Type != ResponseTypeTrade {
			continue
		}
		sme <- response.Data // add message to channel
	}
}

func (t *simpleMovingAverage) process(ch chan []finnhub.FinnRespData) {
	fmt.Println("Start calculating simple moving average")
	t.Window = WindowSize
	t.LastPrices = make(map[string][]float64, 0)
	t.SimpleMovingAverage = make(map[string]float64, 0)

	iteration := 0
	for {
		t.Lock()
		iteration += 1
		data := <-ch

		lastPrices := make(map[string][]float64, 0)
		for _, d := range data {
			if _, ok := lastPrices[d.S]; !ok {
				lastPrices[d.S] = []float64{}
			}
			lastPrices[d.S] = append(lastPrices[d.S], d.P)
			if len(lastPrices[d.S]) > t.Window {
				lastPrices[d.S] = lastPrices[d.S][1:]
			}
		}

		for k, v := range lastPrices {
			avg, err := getRespAverage(v)
			if err != nil {
				log.Print(err, " skipping...")
				errors.Log(err)
				t.Unlock()
				continue
			}
			t.LastPrices[k] = append(t.LastPrices[k], avg)
		}

		if iteration%t.Window == 0 {

			t.calcSimpleMovingAverage()

			for k, v := range t.SimpleMovingAverage {
				fmt.Printf("Simple moving average [%s]: %f \n", k, v)
				out := fmt.Sprintf("SMA [%s]: %f time: %s \n", k, v, time.Now().String())
				err := datastore.WriteFile(out, "sma_records.txt")
				if err != nil {
					log.Print(err, " skipping...")
					errors.Log(fmt.Errorf("%s write data to file: %s \n", time.Now().String(), err.Error()))
					t.Unlock()
					continue
				}
			}
		}
		t.Unlock()
	}
}

func getRespAverage(resp []float64) (float64, error) {
	if len(resp) == 0 {
		return 0, fmt.Errorf("%s get response average: no response data \n", time.Now().String())
	}
	var sum float64
	for _, d := range resp {
		sum += d
	}
	return sum / float64(len(resp)), nil
}

func (t *simpleMovingAverage) calcSimpleMovingAverage() {
	var sum float64
	for k, v := range t.LastPrices {
		for _, p := range v {
			sum += p
		}
		t.SimpleMovingAverage[k] = sum / float64(len(v))
	}
	t.LastPrices = make(map[string][]float64, 0)
}
