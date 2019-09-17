import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';

const baseUrl = '/api/v1/hostInfo/';

@Injectable({
  providedIn: 'root'
})
export class HostInfoService {

  constructor(private http: HttpClient) {
  }

}
