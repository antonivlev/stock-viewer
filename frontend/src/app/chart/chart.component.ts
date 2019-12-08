import { Component, OnInit, OnChanges, Input } from '@angular/core';

@Component({
  selector: 'app-chart',
  templateUrl: './chart.component.html',
  styleUrls: ['./chart.component.css']
})
export class ChartComponent implements OnChanges {
  @Input() data = {};

  ngOnChanges() {
    console.log("app-chart: got this data")
    console.log(this.data);
    this.plotData(this.data);
  }

  plotData(datesMap) {
    let [x, open, high, low, close] = [
        [],
        [],
        [],
        [],
        []
    ];
    Object.keys(datesMap).map(date => {
        x.push(date);
        open.push(datesMap[date]["1. open"]);
        high.push(datesMap[date]["2. high"]);
        low.push(datesMap[date]["3. low"]);
        close.push(datesMap[date]["4. close"]);
    });
    console.log("app-chart: arranged data")

    let trace1 = {
        // data
        x: x,
        open: open,
        high: high,
        low: low,
        close: close,
        // config
        decreasing: { line: { color: '#7F7F7F' } },
        increasing: { line: { color: '#17BECF' } },
        line: { color: 'rgba(31,119,180,1)' },
        type: 'candlestick',
        xaxis: 'x',
        yaxis: 'y'
    };

    let data = [trace1];

    let layout = {
        dragmode: 'zoom',
        margin: { r: 0, t: 0, b: 0, l: 30 },
        showlegend: false,
        xaxis: {
            autorange: true,
            domain: [0, 1],
            title: 'Date',
            type: 'date'
        },
        yaxis: {
            type: 'linear'
        }
    };

    console.log("app-chart: plotly data and layout: ", trace1, layout);
}
}
