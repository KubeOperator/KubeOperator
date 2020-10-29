import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {UserNotificationConfig} from '../message';

@Injectable({
    providedIn: 'root'
})
export class UserSubscribeService {

    baseUrl = '/api/v1/message/subscribe';

    constructor(private http: HttpClient) {
    }

    singleGet(): Observable<UserNotificationConfig[]> {
        return this.http.get<UserNotificationConfig[]>(this.baseUrl);
    }

    singleUpdate(item): Observable<UserNotificationConfig> {
        return this.http.post<UserNotificationConfig>(this.baseUrl, item);
    }
}
