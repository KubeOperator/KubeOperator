import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {HostInfo} from '../host';

const baseUrl = '/api/v1/hostInfo/';

@Injectable({
  providedIn: 'root'
})
export class HostInfoService {

  constructor(private http: HttpClient) {
  }

  loadHostInfo(hostId: string): Observable<HostInfo> {
    return this.http.post<HostInfo>(baseUrl, {'host': hostId});
  }
}
