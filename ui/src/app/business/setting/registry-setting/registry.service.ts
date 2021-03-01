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

    mixedGet(): Observable<Registry> {
        const itemUrl = `${this.baseUrl}`;
        return this.http.get<Registry>(itemUrl);
    }

    getRegistrtyBy(arch): Observable<string> {
        const itemUrl = `${this.baseUrl}/${arch}`;
        return this.http.get<string>(itemUrl);
    }

    CreateRegistry(): Observable<Registry> {
        const itemUrl = `${this.baseUrl}`;
        return this.http.get<Registry>(itemUrl);
    }
}
