import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {Registry} from './registry';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Page} from '../../../shared/class/Page';

@Injectable({
    providedIn: 'root'
})
export class RegistryService extends BaseModelService<Registry> {
    baseUrl = '/api/v1/settings/registry';
    // items: T[] = [];

    constructor(http: HttpClient) {
        super(http);
    }

    mixedGet(page, size): Observable<Page<Registry>> {
        const itemUrl = `${this.baseUrl}?pageNum=${page}&pageSize=${size}`;
        return this.http.get<Page<Registry>>(itemUrl);
    }

    updateRegistryBy(arch: string, item: Registry): Observable<Registry> {
        const url = `${this.baseUrl}/${arch}`;
        return this.http.patch<Registry>(url, item);
    }
}
