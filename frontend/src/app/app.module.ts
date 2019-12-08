import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
// needed for ngModel
import { FormsModule } from '@angular/forms';

import { AppComponent } from './app.component';
import { StockEntryComponent } from './stock-entry/stock-entry.component';
import { ChartComponent } from './chart/chart.component';
import { ErrorDisplayComponent } from './error-display/error-display.component';

@NgModule({
  declarations: [
    AppComponent,
    StockEntryComponent,
    ChartComponent,
    ErrorDisplayComponent
  ],
  imports: [
    BrowserModule,
    FormsModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
