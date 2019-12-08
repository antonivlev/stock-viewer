package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/antonivlev/stock-viewer/apihelpers"
	"github.com/antonivlev/stock-viewer/database"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	errDb := database.SetupDatabase("config.json")
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

func getStockData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println(r.URL.String())
	// parse stock symbol
	keys, ok := r.URL.Query()["symbol"]
	if !ok {
		apihelpers.WriteError(w, "No stock symbol supplied", nil)
		return
	}
	stock := keys[0]

	// check cache
	if database.IsStockCached(stock) {
		log.Println("	getting stock data from cache")
		stockDataBytes := database.GetCachedStockData(stock)
		w.Header().Set("Content-Type", "application/json")
		w.Write(stockDataBytes)
		return
	}
	log.Println("	fetching stock data from alpha vantage")

	// make request to alpha vantage
	resp, errGet := http.Get("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=" + stock + "&apikey=88ABYOD45M3WPBO4")
	if errGet != nil {
		apihelpers.WriteError(w, "Error accessing alpha vantage api", errGet)
		return
	}
	defer resp.Body.Close()

	// parse response
	var apiResponse map[string]interface{}
	errDecode := json.NewDecoder(resp.Body).Decode(&apiResponse)
	// if could not parse response
	if errDecode != nil {
		apihelpers.WriteError(w, "Error parsing api response", errDecode)
		return
	}
	datesMap, ok := apiResponse["Time Series (Daily)"]
	if !ok {
		// response parsed, but it doesnt have the data in it
		apihelpers.WriteErrorResponse(w, apiResponse)
		return
	}

	// otherwise return time series data
	stockDataBytes, errMarshal := json.Marshal(datesMap)
	if errMarshal != nil {
		apihelpers.WriteError(w, "Error encoding stock data", errMarshal)
		return
	}
	// save stock data to db here
	log.Println("	saving stock data from alpha vantage")
	database.SaveStockData(stock, stockDataBytes)
	w.Header().Set("Content-Type", "application/json")
	w.Write(stockDataBytes)
}

func saveSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println(r.URL.String())

	// TODO: database.Search should probably be private
	var search database.Search
	errDecode := json.NewDecoder(r.Body).Decode(&search)
	// if could not parse response
	if errDecode != nil {
		apihelpers.WriteError(w, "Error parsing request body", errDecode)
		return
	}

	errSave := database.SaveSearch(search)
	if errSave != nil {
		apihelpers.WriteError(w, "Error saving to database", errSave)
	}
}

func getSearches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println(r.URL.String())

	searches, errRead := database.GetSearches()
	if errRead != nil {
		apihelpers.WriteError(w, "Error reading from database", errRead)
		return
	}

	searchesBytes, errMarshal := json.Marshal(searches)
	if errMarshal != nil {
		apihelpers.WriteError(w, "Error encoding data from db", errMarshal)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(searchesBytes)
}

// Error handling thought: this function will only be called with a correct "symbol" parameter; because it is called programmatically
// from the successful searches table. No need to handle bad parameter? Can this guarantee be explicit?
func getLatestStockData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println(r.URL.String())
	// parse stock symbol
	keys, ok := r.URL.Query()["symbol"]
	if !ok {
		apihelpers.WriteError(w, "No stock symbol supplied", nil)
		return
	}

	// make request to alpha vantage
	resp, errGet := http.Get("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" + keys[0] + "&apikey=88ABYOD45M3WPBO4")
	if errGet != nil {
		apihelpers.WriteError(w, "", errGet)
		return
	}
	defer resp.Body.Close()
	// parse response
	var apiResponse map[string]interface{}
	errDecode := json.NewDecoder(resp.Body).Decode(&apiResponse)
	// if could not parse response
	if errDecode != nil {
		apihelpers.WriteError(w, "Could not parse response from alpha vantage", errDecode)
		return
	}

	stockData, ok := apiResponse["Global Quote"]
	if !ok {
		apihelpers.WriteErrorResponse(w, apiResponse)
		return
	}

	stockDataBytes, errMarshal := json.Marshal(stockData)
	if errMarshal != nil {
		apihelpers.WriteError(w, "Error encoding stock data", errMarshal)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(stockDataBytes)
}
