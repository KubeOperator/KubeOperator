import {Injectable} from '@angular/core';
import {BaseModelService} from '../../shared/class/BaseModelService';
import {System} from './system/system';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class SystemService extends BaseModelService<System> {

    baseUrl = '/api/v1/settings';

    constructor(http: HttpClient) {
        super(http);
    }

    singleGet(): Observable<System> {
        const itemUrl = `${this.baseUrl}`;
        return this.http.get<System>(itemUrl);
    }
    
    getRegistry() : any {
        const url = '/api/v1/settings/registry';
        return this.http.get<any>(url);
    }

    checkBy(type, item): Observable<System> {
        const itemUrl = `${this.baseUrl}/check/` + type;
        return this.http.post<System>(itemUrl, item);
    }

    getByTab(tab): Observable<System> {
        const itemUrl = `${this.baseUrl}`;
        return this.http.get<System>(itemUrl + '/' + tab);
    }
}
