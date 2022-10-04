package main

import (
	"testing"
)

var validFinnResponse = finnResp{
	Type: "trade",
	Data: []finnRespData{
		{
			S: "BINANCE:BTCUSDT",
			P: 19883.28,
			T: 1664875388916,
			V: 0.1076,
		},
		{
			S: "BINANCE:BTCUSDT",
			P: 19883.21,
			T: 1664875388916,
			V: 0.00768,
		},
		{
			S: "BINANCE:BTCUSDT",
			P: 19883.2,
			T: 1664875388916,
			V: 0.00542,
		},
	},
}

var invalidFinnResponse = finnResp{
	Type: "trade",
	Data: []finnRespData{},
}

func TestGetRespAverageValid(t *testing.T) {
	_, err := getRespAverage(validFinnResponse)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

}

func TestGetRespAverageInvalidNoData(t *testing.T) {
	avg, err := getRespAverage(invalidFinnResponse)
	if avg != 0 {
		t.Errorf("Error: %s", err)
	}

}
