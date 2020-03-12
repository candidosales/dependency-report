import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { BrowserModule } from '@angular/platform-browser';
import { DataService } from './providers/data.service';
import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { PieChartComponent } from './google-charts/pie-chart.component';
import {MatTooltipModule} from '@angular/material/tooltip';


@NgModule({
  declarations: [
    AppComponent,
    PieChartComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MatTooltipModule
  ],
  providers: [DataService],
  bootstrap: [AppComponent]
})
export class AppModule { }
