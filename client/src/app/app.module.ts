import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';

import { AppComponent } from './app.component';
import { FrappeDirective } from './frappe/frappe.directive';
import { DataService } from './providers/data.service';

@NgModule({
  declarations: [
    AppComponent,
    FrappeDirective
  ],
  imports: [
    BrowserModule,
    HttpClientModule
  ],
  providers: [DataService],
  bootstrap: [AppComponent]
})
export class AppModule { }
