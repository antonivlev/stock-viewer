import { Component, OnInit, Input } from '@angular/core';

@Component({
  selector: 'app-stock-entry',
  templateUrl: './stock-entry.component.html',
  styleUrls: ['./stock-entry.component.css']
})
export class StockEntryComponent implements OnInit {
  @Input() currentSearch: string;

  getSearchData() {
    console.log("fetching search data for: ", this.currentSearch);
    //this.stockDataService.fetchData(this.currentSearch);
    // responsibility ends here
    // service fires event; either "successfult-search" (with data) or "fetch-error" (with error). other components react.
  }

  ngOnInit() {}
}