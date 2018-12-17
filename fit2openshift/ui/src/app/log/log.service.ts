import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Log} from './log';

const baseUrl = '/api/v1/cluster/{clusterId}/log';

@Injectable()
export class LogService {

  constructor(private http: HttpClient) {
  }

  getLogs(clusterId): Observable<Log[]> {
    return this.http.get<Log[]>(`${baseUrl.replace('{clusterId}', clusterId)}`);
  }
}
