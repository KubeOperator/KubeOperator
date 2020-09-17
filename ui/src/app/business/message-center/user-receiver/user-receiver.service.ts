import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {UserReceiver} from '../message';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class UserReceiverService extends BaseModelService<UserReceiver> {

    baseUrl = '/api/v1/message/receiver';

    constructor(http: HttpClient) {
        super(http);
    }

    singleGet(userName): Observable<UserReceiver> {
        const itemUrl = `${this.baseUrl}?userName=${userName}`;
        return this.http.get<UserReceiver>(itemUrl);
    }

    singleUpdate(item): Observable<UserReceiver> {
        const itemUrl = `${this.baseUrl}`;
        return this.http.post<UserReceiver>(itemUrl, item);
    }
}
