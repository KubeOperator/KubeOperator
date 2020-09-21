import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {UserNotificationConfig, UserReceiver} from '../message';

@Injectable({
    providedIn: 'root'
})
export class UserSubscribeService {

    baseUrl = '/api/v1/message/subscribe';

    constructor(private http: HttpClient) {
    }

    singleGet(userName): Observable<UserNotificationConfig[]> {
        const itemUrl = `${this.baseUrl}?userName=${userName}`;
        return this.http.get<UserNotificationConfig[]>(itemUrl);
    }

    singleUpdate(item): Observable<UserNotificationConfig> {
        const itemUrl = `${this.baseUrl}`;
        return this.http.post<UserNotificationConfig>(itemUrl, item);
    }
}
