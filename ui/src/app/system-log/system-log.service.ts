import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {SystemLog, SystemLogPager} from './system-log';

@Injectable({
  providedIn: 'root'
})
export class SystemLogService {

  constructor(private http: HttpClient) {
  }

  private baseUrl = '/api/v1/log/';

  searchLog(params): Observable<SystemLogPager> {
    return this.http.post<SystemLogPager>(this.baseUrl, params);
  }
}
