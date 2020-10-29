import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import {F5CreateRequest} from './f5';

@Injectable({
  providedIn: 'root'
})
export class F5Service {
  baseUrl = '/api/v1/f5';

  getItems(clusterName: string): Observable<F5CreateRequest[]>{
    const itemUrl = `${this.baseUrl}/${clusterName}`;
    return this.http.get<F5CreateRequest[]>(itemUrl);
  }
  constructor(private  http: HttpClient) {
  }
}
