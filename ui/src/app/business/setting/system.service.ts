import {Injectable} from '@angular/core';
import {BaseModelService} from '../../shared/class/BaseModelService';
import {System} from './system/system';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class SystemService extends BaseModelService<System> {

    baseUrl = '/api/v1/systemSettings';

    constructor(http: HttpClient) {
        super(http);
    }

    singleGet(): Observable<System> {
        const itemUrl = `${this.baseUrl}`;
        return this.http.get<System>(itemUrl);
    }

    getIp(): Observable<string> {
        const itemUrl = `${this.baseUrl}/ip`;
        return this.http.get<string>(itemUrl);
    }
}