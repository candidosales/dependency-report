import { Data } from '../interface/data.interface';
import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class DataService {

  constructor(private httpClient: HttpClient) { }


  public getDataInServer(): Observable<Data> {
    return this.httpClient.get<Data>('/assets/config/data-test.json');
  }
}
