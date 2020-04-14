import { Data } from '../interface/data.interface';
import { environment } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class DataService {
  apiUrl = '';

  constructor(private httpClient: HttpClient) {
    this.apiUrl = environment.apiUrl;
  }

  public getCacheData(): Observable<Data> {
    return this.httpClient.get<Data>('/assets/config/data.json');
  }

  public generateReport(): Observable<Data> {
    return this.httpClient.get<Data>(`${this.apiUrl}/generate-report`);
  }

  public ping() {
    return this.httpClient.get<any>(`${this.apiUrl}/ping`);
  }
}
