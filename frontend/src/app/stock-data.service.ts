import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class StockDataService {

  constructor() { }

  fetchStockData(symbol:string) {
    console.log("stock-data-service: fetching in service");
    // fetch here, on success dispatch "record-search", "plot-data"
    // on failure dispatch "server-error"
    return fetch("http://localhost:3000/api/get-stock-data?symbol=" + symbol)
      .then(res => {
          if (res.status === 200) {
              // all good
              return res.json().then(stockData => ({ok: true, data: stockData}));
          } else {
              // server returned error
              return res.text().then(errorText => ({ok: false, data: errorText}));
          }
      })
    }

  saveSearch(symbol:string, time:string) {
    //
  }
}
