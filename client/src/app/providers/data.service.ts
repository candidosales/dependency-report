import { Data } from '../interface/data.interface';
import { environment } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { of } from 'rxjs';
import { shareReplay, map, catchError } from 'rxjs/operators';
import { MatSnackBar } from '@angular/material/snack-bar';

export interface PingData {
  ok: boolean;
}

@Injectable({
  providedIn: 'root'
})
export class DataService {

  cacheData$ = this.httpClient.get<Data>('/assets/config/data.json').pipe(
    map(value => {
      this.prepareData(value);
      return value;
    }),
    // shareReplay(1)
  );

  data$ = this.httpClient.get<Data>(`${environment.apiUrl}/generate-report`).pipe(
    map(value => {
      this.prepareData(value);
      return value;
    }),
    shareReplay(1)
  );

  serverIsOn$ = this.httpClient.get<PingData>(`${environment.apiUrl}/ping`).pipe(
    map(value => {
      return value.ok;
    }),
    catchError(err => {
      console.error(err);
      this.snackBar.open('Sorry, but the server is not working', 'CLOSE', {
        duration: 10000
      });
      return of(false);
    }),
    shareReplay(1)
  );

  constructor(private httpClient: HttpClient, private snackBar: MatSnackBar) {
  }

  prepareData(data: Data) {
    console.log('data', data);
    data.graphData.projectsByFilters?.unshift(['Filter', 'Version']);
    data.graphData.componentsByFilters?.unshift(['Filter', 'Version']);
    Object.keys(data.dependenciesByVersions).filter(key => {
      const quantity = Object.keys(data.dependenciesByVersions[key]).length;
      data.dependenciesByVersions[key].type = this.getDependencyClassification(quantity);
    });
  }

  getDependencyClassification(quantity: number): string {
    if (quantity <= 2) {
      return 'good';
    }
    if (quantity > 2 && quantity <= 5) {
      return 'warning';
    }
    if (quantity > 5 && quantity <= 10) {
      return 'bad';
    }
    if (quantity > 10) {
      return 'terrible';
    }
    return '';
  }
}
