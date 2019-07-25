import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Version} from './header/version';

@Injectable({
  providedIn: 'root'
})
export class BaseService {
  baseUrl = '/api/v1/version/';

  constructor(private http: HttpClient) {
  }

  getVersion(): Observable<Version> {
    return this.http.get<Version>(this.baseUrl);
  }
}
