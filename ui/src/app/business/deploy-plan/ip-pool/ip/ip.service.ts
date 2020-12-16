import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../../shared/class/BaseModelService';
import {Ip} from './ip';
import {HttpClient} from '@angular/common/http';
import {Page} from '../../../../shared/class/Page';
import {Observable} from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class IpService extends BaseModelService<Ip> {

    baseUrl = '/api/v1/ippools/{name}/ips';

    constructor(http: HttpClient) {
        super(http);
    }

    page(page, size, ipPoolName?: string): Observable<Page<Ip>> {
        let url = this.baseUrl;
        if (ipPoolName) {
            url = this.baseUrl.replace('{name}', ipPoolName);
        }
        const pageUrl = `${url}?pageNum=${page}&pageSize=${size}`;
        return this.http.get<Page<Ip>>(pageUrl);
    }
}
