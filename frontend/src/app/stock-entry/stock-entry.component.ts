import { Component, OnInit, Input, Output } from '@angular/core';
import { StockDataService } from '../stock-data.service';
import { EventEmitter } from '@angular/core';

@Component({
  selector: 'app-stock-entry',
  templateUrl: './stock-entry.component.html',
  styleUrls: ['./stock-entry.component.css']
})
export class StockEntryComponent implements OnInit {
  @Input() currentSearch: string;
  @Output() newData = new EventEmitter();
  @Output() newError = new EventEmitter();

  constructor(private stockDataService: StockDataService) {}

  getSearchData() : void {
    console.log("app-stock-entry: fetching search data for: ", this.currentSearch);
    this.stockDataService.fetchStockData(this.currentSearch)
      .then(result => {
        if (result.ok) {
          this.newData.emit(result.data)
          this.newError.emit({})
        } else {
          this.newError.emit(result.data)
        }
      })    
  }

  ngOnInit() {}
}