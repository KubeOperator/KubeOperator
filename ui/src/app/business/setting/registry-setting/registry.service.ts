import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {Registry} from './registry';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class RegistryService extends BaseModelService<Registry> {
    baseUrl = '/api/v1/settings/registry';

    constructor(http: HttpClient) {
        super(http);
    }

    // mixedGet(): Observable<Registry> {
    //     const itemUrl = `${this.baseUrl}`;
    //     return this.http.get<Registry>(itemUrl);
    // }
    //
    updateRegistryBy(arch: string, item: Registry): Observable<Registry> {
        const url = `${this.baseUrl}/${arch}`;
        return this.http.patch<Registry>(url, item);
    }

    // //
    // // update(name: string, item: BaseRequest, ipPoolName?: string): Observable<Ip> {
    // //     let url = this.baseUrl;
    // //     if (ipPoolName) {
    // //         url = this.baseUrl.replace('{name}', ipPoolName);
    // //     }
    // //     return this.http.patch<Ip>(url, item);
    // // }
    //
    // CreateRegistry(): Observable<Registry> {
    //     const itemUrl = `${this.baseUrl}`;
    //     return this.http.get<Registry>(itemUrl);
    // }
}
