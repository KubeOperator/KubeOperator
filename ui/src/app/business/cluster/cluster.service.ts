import {Injectable} from '@angular/core';
import {BaseModelService} from '../../shared/class/BaseModelService';
import {HttpClient} from '@angular/common/http';
import {
    Cluster, ClusterHealthCheck,
    CLusterImportRequest, ClusterRecoverItem,
    ClusterSecret,
    ClusterStatus,
    ClusterUpgradeRequest,
    InitClusterResponse
} from './cluster';
import {Observable} from 'rxjs';
import {Page} from '../../shared/class/Page';
import {ClusterLog} from './cluster-detail/log/log';
import {Batch} from "../../shared/class/Batch";

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

    secret(clusterName: string): Observable<ClusterSecret> {
        return this.http.get<ClusterSecret>(`${this.baseUrl}/secret/${clusterName}`);
    }

    log(clusterName: string): Observable<ClusterLog[]> {
        return this.http.get<ClusterLog[]>(`${this.baseUrl}/log/${clusterName}`);
    }


    pageBy(page, size, projectName): Observable<Page<Cluster>> {
        const pageUrl = `${this.baseUrl}?pageNum=${page}&pageSize=${size}&projectName=${projectName}`;
        return this.http.get<Page<Cluster>>(pageUrl);
    }

    upgrade(clusterName: string, version: string): Observable<any> {
        const req = new ClusterUpgradeRequest();
        req.clusterName = clusterName;
        req.version = version;
        return this.http.post<any>(`${this.baseUrl}/upgrade/`, req);
    }

    healthCheck(clusterName: string): Observable<ClusterHealthCheck> {
        return this.http.get<ClusterHealthCheck>(`${this.baseUrl}/health/${clusterName}`);
    }

    recover(clusterName: string): Observable<ClusterRecoverItem[]> {
        return this.http.post<ClusterRecoverItem[]>(`${this.baseUrl}/recover/${clusterName}`, {});
    }


    batchDelete(method: string, items: Cluster[], force: boolean, projectName?: string): Observable<any> {
        const options = {};
        if (projectName) {
            options['headers'] = {
                project: encodeURI(projectName)
            };
        }
        const batchUrl = `${this.baseUrl}/batch/?force=true`;
        const b = new Batch<Cluster>(method, items);
        return this.http.post(batchUrl, b, options);
    }
}



