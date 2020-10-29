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

    singleGet(): Observable<UserReceiver> {
        return this.http.get<UserReceiver>(this.baseUrl);
    }

    singleUpdate(item): Observable<UserReceiver> {
        return this.http.post<UserReceiver>(this.baseUrl, item);
    }
}
