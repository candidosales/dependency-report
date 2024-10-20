import { Data } from '../interface/data.interface';
import { environment } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { of } from 'rxjs';
import { shareReplay, map, catchError } from 'rxjs/operators';
import { MatLegacySnackBar as MatSnackBar } from '@angular/material/legacy-snack-bar';

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
    shareReplay(1)
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
    data.graphData.projectsByFilters?.unshift(['Filter', 'Version']);
    data.graphData.componentsByFilters?.unshift(['Filter', 'Version']);
  }
}
