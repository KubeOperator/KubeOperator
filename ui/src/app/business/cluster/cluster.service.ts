import {Injectable} from '@angular/core';
import {BaseModelService} from '../../shared/class/BaseModelService';
import {HttpClient} from '@angular/common/http';
import {Cluster, CLusterImportRequest, ClusterStatus, InitClusterResponse} from './cluster';
import {Observable} from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class ClusterService extends BaseModelService<Cluster> {

    baseUrl = '/api/v1/clusters';

    constructor(http: HttpClient) {
        super(http);
    }

    status(clusterName: string): Observable<ClusterStatus> {
        return this.http.get<ClusterStatus>(`${this.baseUrl}/status/${clusterName}`);
    }

    init(clusterName: string): Observable<InitClusterResponse> {
        return this.http.post<InitClusterResponse>(`${this.baseUrl}/init/${clusterName}/`, {});
    }

    import(item: CLusterImportRequest): Observable<any> {
        return this.http.post<any>(`${this.baseUrl}/import/`, item);
    }
}
