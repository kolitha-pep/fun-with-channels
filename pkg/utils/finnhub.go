package utils

import (
	"encoding/json"

	"github.com/kolitha-pep/fun-with-channels/app/pkg/finnhub"
)

func FinnRespToStruct(chanData interface{}) finnhub.FinnResp {
	m, err := json.Marshal(chanData)
	if err != nil {
		panic(err)
	}

	var resp finnhub.FinnResp
	err = json.Unmarshal(m, &resp)
	if err != nil {
		panic(err)
	}
	return resp
}
