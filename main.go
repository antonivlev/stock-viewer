package main

import (
	"encoding/json"
	"fmt"
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

	fmt.Println("Serving on :3000")
	http.ListenAndServe(":3000", nil)
}

func getStockData(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())
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

	// parse response
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
	// extract stock data from response, save back to json
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

func saveSearch(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())

	var search database.Search
	errDecode := json.NewDecoder(r.Body).Decode(&search)
	// if could not parse response
	if errDecode != nil {
		fmt.Fprintf(w, "Error parsing request body:\n %s", errDecode.Error())
		return
	}

	errSave := database.SaveSearch(search)
	if errSave != nil {
		fmt.Fprintf(w, "Error saving to database: \n %v", errSave.Error())
	}
}

func getSearches(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())

	searches, errRead := database.GetSearches()
	if errRead != nil {
		fmt.Fprintf(w, "Error reading from database: \n %v", errRead.Error())
		return
	}

	// fetch and add current stock data to every search in db
	// TODO: change, very inefficient. Use getLatestStockData from client
	// type searchWithData struct {
	// 	database.Search
	// 	latestData
	// }
	// var infoSearches []searchWithData
	// for _, s := range searches {
	// 	latest, _ := getLatestStockData(s.Stock)
	// 	infoSearches = append(
	// 		infoSearches,
	// 		searchWithData{s, latest},
	// 	)
	// }

	searchesBytes, errMarshal := json.Marshal(searches)
	if errMarshal != nil {
		fmt.Fprintf(w, "Error encoding data from db:\n %s", errMarshal.Error())
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

func getLatestStockData(stock string) (latestData, error) {
	// make request to alpha vantage
	resp, errGet := http.Get("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" + stock + "&apikey=88ABYOD45M3WPBO4")
	if errGet != nil {
		return latestData{}, errGet
	}
	defer resp.Body.Close()
	// parse response
	var stockData struct {
		GlobalQuote latestData `json:"Global Quote"`
	}
	//var stockData map[string]interface{}
	errDecode := json.NewDecoder(resp.Body).Decode(&stockData)
	// if could not parse response
	if errDecode != nil {
		return latestData{}, errDecode
	}
	// TODO: error in reponse not handled

	return stockData.GlobalQuote, nil
}
