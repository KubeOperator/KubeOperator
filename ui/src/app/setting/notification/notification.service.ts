import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Settings} from '../setting';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class NotificationService {
  baseUrl = '/api/v1/notification/';

  constructor(private http: HttpClient) {
  }

  emailCheck(email: Settings): Observable<Settings> {
    const url = this.baseUrl + 'email/check/';
    return this.http.post<Settings>(url, email);
  }

  workWeixinCheck(workWeixin: Settings): Observable<Settings> {
    const url = this.baseUrl + 'workWeixin/check/';
    return this.http.post<Settings>(url, workWeixin);
  }
}
