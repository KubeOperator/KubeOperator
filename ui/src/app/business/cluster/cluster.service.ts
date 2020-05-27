import {Injectable} from '@angular/core';
import {BaseModelService} from '../../shared/class/BaseModelService';
import {HttpClient} from '@angular/common/http';
import {Cluster, ClusterStatusResponse, InitClusterResponse, Status} from './cluster';
import {Observable} from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class ClusterService extends BaseModelService<Cluster> {

    baseUrl = '/api/v1/clusters';

    constructor(http: HttpClient) {
        super(http);
    }

    status(clusterName: string): Observable<ClusterStatusResponse> {
        return this.http.get<ClusterStatusResponse>(`${this.baseUrl}/${clusterName}/status/`);
    }

    init(clusterName: string): Observable<InitClusterResponse> {
        return this.http.post<InitClusterResponse>(`${this.baseUrl}/init/${clusterName}/`, {});
    }
}
