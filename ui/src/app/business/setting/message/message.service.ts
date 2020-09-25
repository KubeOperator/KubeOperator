import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {System} from '../system/system';

@Injectable({
    providedIn: 'root'
})
export class MessageService {

    baseUrl = '/api/v1/message/setting';

    constructor(private http: HttpClient) {
    }

    getByTab(tabName): Observable<System> {
        const itemUrl = `${this.baseUrl}/${tabName}`;
        return this.http.get<System>(itemUrl);
    }

    postByTab(tabName, item): Observable<System> {
        const itemUrl = `${this.baseUrl}/${tabName}`;
        return this.http.post<System>(itemUrl, item);
    }

    postCheckByTab(tabName, item): Observable<System> {
        const itemUrl = `${this.baseUrl}/check/${tabName}`;
        return this.http.post<System>(itemUrl, item);
    }
}
