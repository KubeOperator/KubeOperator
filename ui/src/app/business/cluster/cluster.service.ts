import {Injectable} from '@angular/core';
import {BaseModelService} from '../../shared/class/BaseModelService';
import {HttpClient} from '@angular/common/http';
import {
    Cluster,
    CLusterImportRequest,
    ClusterSecret,
    ClusterStatus,
    ClusterUpgradeRequest,
    InitClusterResponse
} from './cluster';
import {Observable} from 'rxjs';
import {Page} from '../../shared/class/Page';

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
}
