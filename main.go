package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/antonivlev/stock-viewer/database"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	errDb := database.SetupDatabase()
	if errDb != nil {
		fmt.Printf("Error connecting to database: \n%s\n", errDb.Error())
	}

	// serve the main html file
	http.Handle("/", http.FileServer(http.Dir("./static")))
	// api, frontend calls these for data
	http.HandleFunc("/api/get-stock-data", getStockData)
	http.HandleFunc("/api/save-search", saveSearch)
	http.HandleFunc("/api/get-searches", getSearches)
	http.HandleFunc("/api/get-latest-stock-data", getLatestStockData)

	fmt.Println("Serving on localhost:3000")
	http.ListenAndServe(":3000", nil)
}

// Writes message, error text and Bad Request code. If err is nil, just writes message.
func writeError(w http.ResponseWriter, message string, err error) {
	errString := ""
	if err != nil {
		errString = err.Error()
	}
	http.Error(w, message+"\n\n"+errString, http.StatusBadRequest)
}

func getStockData(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.String())
	// parse stock symbol
	keys, ok := r.URL.Query()["symbol"]
	if !ok {
		writeError(w, "No stock symbol supplied", nil)
		return
	}

	// make request to alpha vantage
	resp, errGet := http.Get("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=" + keys[0] + "&apikey=88ABYOD45M3WPBO4")
	if errGet != nil {
		writeError(w, "Error accessing alpha vantage api", errGet)
		return
	}
	defer resp.Body.Close()

	// parse response
	var stockData map[string]interface{}
	errDecode := json.NewDecoder(resp.Body).Decode(&stockData)
	// if could not parse response
	if errDecode != nil {
		writeError(w, "Error parsing stock data response", errDecode)
		return
	}
	// if response contains error (e.g. wrong stock symbol)
	errAPIText, hasError := stockData["Error Message"].(string)
	if hasError {
		writeError(w, "Alpha Vantage API returned an error. Is your stock symbol valid?"+errAPIText, nil)
		return
	}
	// extract stock data from response, save back to json
	// TODO: might need to handle access error here, check apha vantage api guarantees
	datesMap := stockData["Time Series (Daily)"]
	stockDataBytes, errMarshal := json.Marshal(datesMap)
	if errMarshal != nil {
		writeError(w, "Error encoding stock data", errMarshal)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(stockDataBytes)
}

func saveSearch(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.String())

	var search database.Search
	errDecode := json.NewDecoder(r.Body).Decode(&search)
	// if could not parse response
	if errDecode != nil {
		writeError(w, "Error parsing request body", errDecode)
		return
	}

	errSave := database.SaveSearch(search)
	if errSave != nil {
		writeError(w, "Error saving to database", errSave)
	}
}

func getSearches(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.String())

	searches, errRead := database.GetSearches()
	if errRead != nil {
		writeError(w, "Error reading from database", errRead)
		return
	}

	searchesBytes, errMarshal := json.Marshal(searches)
	if errMarshal != nil {
		writeError(w, "Error encoding data from db", errMarshal)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(searchesBytes)
}

type latestData struct {
	Open   string `json:"02. open"`
	High   string `json:"03. high"`
	Low    string `json:"04. low"`
	Close  string `json:"08. previous close"`
	Volume string `json:"06. volume"`
}

// TODO: very simillar to getStockData, these two should share functions; e.g. response parsing
// Error handling thought: this function will only be called with a correct "symbol" parameter; because it is called programmatically
// from the successful searches table. No need to handle bad parameter? Can this guarantee be explicit?
func getLatestStockData(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.String())
	// parse stock symbol
	keys, ok := r.URL.Query()["symbol"]
	if !ok {
		writeError(w, "No stock symbol supplied", nil)
		return
	}

	// make request to alpha vantage
	resp, errGet := http.Get("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" + keys[0] + "&apikey=88ABYOD45M3WPBO4")
	if errGet != nil {
		writeError(w, "", errGet)
		return
	}
	defer resp.Body.Close()
	// parse response
	var stockData struct {
		GlobalQuote latestData `json:"Global Quote"`
	}
	errDecode := json.NewDecoder(resp.Body).Decode(&stockData)
	// if could not parse response
	if errDecode != nil {
		writeError(w, "", errDecode)
		return
	}
	// TODO: error in reponse not handled
	stockDataBytes, errMarshal := json.Marshal(stockData.GlobalQuote)
	if errMarshal != nil {
		writeError(w, "Error encoding stock data", errMarshal)
		return
	}
	fmt.Printf("%+v\n", stockData.GlobalQuote)
	w.Header().Set("Content-Type", "application/json")
	w.Write(stockDataBytes)
}
