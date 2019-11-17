package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	// api
	http.HandleFunc("/api/get-stock-data", getStockData)

	fmt.Println("Serving on :3000")
	http.ListenAndServe(":3000", nil)
}

func getStockData(w http.ResponseWriter, r *http.Request) {
	// parse stock symbol
	keys, ok := r.URL.Query()["symbol"]
	if !ok {
		fmt.Fprintf(w, "No stock symbol supplied")
		return
	}

	// make request to alpha vantage
	resp, errGet := http.Get("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=" + keys[0] + "&apikey=88ABYOD45M3WPBO4")
	if errGet != nil {
		fmt.Fprintf(w, "Error accessing alpha vantage api:\n %s", errGet.Error())
		return
	}
	defer resp.Body.Close()

	var stockData map[string]interface{}
	errDecode := json.NewDecoder(resp.Body).Decode(&stockData)
	// if could not parse response
	if errDecode != nil {
		fmt.Fprintf(w, "Error parsing stock data response:\n %s", errDecode.Error())
		return
	}
	// if response contains error (e.g. wrong stock symbol)
	errAPI, hasError := stockData["Error Message"].(string)
	if hasError {
		fmt.Fprintf(w, "Alpha Vantage API returned an error. Is your stock symbol valid? \n %s", errAPI)
		return
	}
	// TODO: might need to handle access error here, check apha vantage api guarantees
	datesMap := stockData["Time Series (Daily)"]

	stockDataBytes, errMarshal := json.Marshal(datesMap)
	if errMarshal != nil {
		fmt.Fprintf(w, "Error encoding stock data:\n %s", errMarshal.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(stockDataBytes)
}
