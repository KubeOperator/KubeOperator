import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Dns} from './dns';

@Injectable({
  providedIn: 'root'
})
export class DnsService {

  constructor(private http: HttpClient) {
  }

  getDns(): Observable<Dns> {
    return this.http.get<Dns>('api/v1/dns/');
  }

  updateDns(dns: Dns): Observable<Dns> {
    return this.http.post<Dns>('api/v1/dns/update/', dns);
  }
}
