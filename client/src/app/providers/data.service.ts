import { Data } from '../interface/data.interface';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class DataService {
  apiUrl = '';

  constructor(private httpClient: HttpClient) {
    this.apiUrl = environment.apiUrl;
  }

  public getDataInServer(): Observable<Data> {
    return this.httpClient.get<Data>('/assets/config/data-test.json');
  }

  public generateReport(): Observable<HttpResponse<any>> {
    return this.httpClient.get<HttpResponse<any>>(`${this.apiUrl}/generate-report`);
  }
}
