import { Data } from '../interface/data.interface';
import { environment } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { shareReplay, map } from 'rxjs/operators';

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
    shareReplay(1)
  );

  constructor(private httpClient: HttpClient) {
  }

  prepareData(value: Data) {
    value.graphData.projectsByFilters?.unshift(['Filter', 'Version']);
    value.graphData.componentsByFilters?.unshift(['Filter', 'Version']);
  }
}
