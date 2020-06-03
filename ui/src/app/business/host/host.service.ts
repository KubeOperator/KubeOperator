import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BaseModelService} from '../../shared/class/BaseModelService';
import {Observable} from 'rxjs';
import {Host, HostCreateRequest} from './host';

@Injectable({
    providedIn: 'root'
})
export class HostService extends BaseModelService<any> {

    baseUrl = '/api/v1/hosts';

    constructor(http: HttpClient) {
        super(http);
    }

    sync(name: string, item: HostCreateRequest): Observable<Host> {
        const itemUrl = `${this.baseUrl}/${name}/sync/`;
        return this.http.post<Host>(itemUrl, item);
    }
}
