package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	FinnAPIKey        = "ccpcv5qad3i91ts8jhk0ccpcv5qad3i91ts8jhkg"
	ResponseTypeTrade = "trade"
	WindowSize        = 10
)

type finnResp struct {
	Type string         `json:"type"`
	Data []finnRespData `json:"data"`
}

type finnRespData struct {
	S string  `json:"s"`
	P float64 `json:"p"`
	T int64   `json:"t"`
	V float64 `json:"v"`
}

type SimpleMovingAverage struct {
	LastPrices          []float64
	Window              int
	SimpleMovingAverage float64
}

func main() {

	// connect to finnhub websocket
	w, _, err := websocket.DefaultDialer.Dial("wss://ws.finnhub.io?token="+FinnAPIKey, nil)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	symbols := []string{"BINANCE:BTCUSDT"}
	for _, s := range symbols {
		msg, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})
		w.WriteMessage(websocket.TextMessage, msg)
	}

	// create a channel and pass it to process function
	sme := make(chan interface{})
	defer close(sme)
	go process(sme)
	var msg interface{}

	for {
		err := w.ReadJSON(&msg)
		if err != nil {
			panic(err)
		}
		sme <- msg
	}
}

func process(ch chan interface{}) {
	fmt.Println("Start calculating simple moving average")
	sma := SimpleMovingAverage{Window: WindowSize}
	var lock sync.Mutex

	iteration := 0
	for {
		lock.Lock()
		iteration += 1
		chanData := <-ch

		response := finnRespToStruct(chanData)
		respAverage, err := getRespAverage(response)
		if err != nil {
			log.Print(err, " skipping...")
			lock.Unlock()
			continue
		}

		sma.LastPrices = append(sma.LastPrices, respAverage)

		if response.Type == ResponseTypeTrade {
			sma.calcSimpleMovingAverage()
			if iteration%sma.Window == 0 {
				writeFile(fmt.Sprintf("SMA: %f time: %s \n", sma.SimpleMovingAverage, time.Now().String()))
				fmt.Println("SMA: ", sma.SimpleMovingAverage)
			}
		}
		lock.Unlock()
	}
}

func finnRespToStruct(chanData interface{}) finnResp {
	m, err := json.Marshal(chanData)
	if err != nil {
		panic(err)
	}

	var resp finnResp
	err = json.Unmarshal(m, &resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func getRespAverage(resp finnResp) (float64, error) {
	if len(resp.Data) == 0 {
		return 0, fmt.Errorf("no response data")
	}

	var sum float64
	for _, d := range resp.Data {
		sum += d.P
	}
	return sum / float64(len(resp.Data)), nil
}

func (sma *SimpleMovingAverage) calcSimpleMovingAverage() float64 {
	var sum float64
	for _, p := range sma.LastPrices {
		sum += p
	}
	sma.SimpleMovingAverage = sum / float64(len(sma.LastPrices))
	sma.LastPrices = nil

	return sma.SimpleMovingAverage
}

func writeFile(string string) {
	// open file for writing simple moving average values
	f, err := os.OpenFile("sma_records.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(string)
	w.Flush()
}
