import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Host} from './host';


const baseUrl = '/api/v1/host/';

@Injectable({
  providedIn: 'root'
})
export class HostService {

  constructor(private http: HttpClient) {

  }

  listHosts(): Observable<Host[]> {
    return this.http.get<Host[]>(baseUrl);
  }

  createHost(host: Host): Observable<Host> {
    return this.http.post<Host>(baseUrl, host);
  }

  deleteHost(hostId: string): Observable<any> {
    return this.http.delete<any>(baseUrl + hostId + '/');
  }


}
