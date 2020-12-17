import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../../shared/class/BaseModelService';
import {Ip} from './ip';
import {HttpClient} from '@angular/common/http';
import {Page} from '../../../../shared/class/Page';
import {Observable} from 'rxjs';
import {Batch} from '../../../../shared/class/Batch';
import {BaseRequest} from '../../../../shared/class/BaseModel';

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

    batch(method: string, items: Ip[], ipPoolName?: string): Observable<any> {
        let batchUrl = this.baseUrl + '/batch/';
        if (ipPoolName) {
            batchUrl = batchUrl.replace('{name}', ipPoolName);
        }
        const b = new Batch<Ip>(method, items);
        return this.http.post(batchUrl, b);
    }

    create(item: BaseRequest, ipPoolName?: string): Observable<Ip> {
        let url = this.baseUrl;
        if (ipPoolName) {
            url = this.baseUrl.replace('{name}', ipPoolName);
        }
        return this.http.post<Ip>(url, item);
    }
}
