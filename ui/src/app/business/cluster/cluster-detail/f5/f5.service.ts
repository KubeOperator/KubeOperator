import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {F5} from './f5';

@Injectable({
    providedIn: 'root'
})
export class F5Service {
    baseUrl = '/api/v1/clusters/f5';

    getItems(clusterName: string): Observable<F5> {
        const itemUrl = `${this.baseUrl}/${clusterName}`;
        return this.http.get<F5>(itemUrl);
    }

    create(item: F5): Observable<F5> {
        const url = this.baseUrl;
        return this.http.post<F5>(url, item);
    }

    update(item: F5): Observable<F5> {
        const url = this.baseUrl;
        return this.http.patch<F5>(url, item);
    }

    constructor(private  http: HttpClient) {
    }
}
