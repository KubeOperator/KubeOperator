import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {HttpClient} from '@angular/common/http';
import {CloudZoneRequest, Zone} from './zone';
import {Observable} from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class ZoneService extends BaseModelService<Zone> {

    baseUrl = '/api/v1/zones';

    constructor(http: HttpClient) {
        super(http);
    }

    listClusters(item: CloudZoneRequest): Observable<any> {
        const itemUrl = `${this.baseUrl}/clusters/`;
        return this.http.post<any>(itemUrl, item);
    }

    listTemplates(item: CloudZoneRequest): Observable<any> {
        const itemUrl = `${this.baseUrl}/templates/`;
        return this.http.post<any>(itemUrl, item);
    }

    listByRegionName(regionName: string): Observable<any> {
        const itemUrl = `${this.baseUrl}/list/` + regionName + '/';
        return this.http.get<any>(itemUrl);
    }

    listDatastores(item: CloudZoneRequest): Observable<any> {
        const itemUrl = `${this.baseUrl}/datastores/`;
        return this.http.post<any>(itemUrl, item);
    }

}
