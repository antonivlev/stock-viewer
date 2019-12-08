import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'Stock Data Viewer';
  chartData = {};
  fetchError: string;

  fillInError(err) {
    this.fetchError = err;
    console.log("app-root: my child (stock-entry) updated my error");
  }

  drawData(data) {
    this.chartData = data;
    console.log("app-root: my child (stock-entry) updated my data");
  }
}
