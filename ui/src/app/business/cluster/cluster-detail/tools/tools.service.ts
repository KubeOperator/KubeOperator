import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';
import {ClusterTool} from './tools';
import {HttpClient} from '@angular/common/http';

@Injectable({
    providedIn: 'root'
})
export class ToolsService {

    constructor(private http: HttpClient) {
    }

    baseUrl = '/api/v1/clusters/tool/{operation}/{cluster_name}';

    list(clusterName: string): Observable<ClusterTool[]> {
        return this.http.get<ClusterTool[]>(this.baseUrl.replace('/{operation}', '').replace('{cluster_name}', clusterName));
    }

    enable(clusterName: string, item: ClusterTool): Observable<ClusterTool> {
        return this.http.post<ClusterTool>(this.baseUrl.replace('{operation}', 'enable').replace('{cluster_name}', clusterName), item);
    }

    upgrade(clusterName: string, item: ClusterTool): Observable<ClusterTool> {
        return this.http.post<ClusterTool>(this.baseUrl.replace('{operation}', 'upgrade').replace('{cluster_name}', clusterName), item);
    }

    disable(clusterName: string, item: ClusterTool): Observable<any> {
        return this.http.post<any>(this.baseUrl.replace('{operation}', 'disable').replace('{cluster_name}', clusterName), item);
    }

    getNodeport(clusterName: string, name: string): Observable<any> {
        return this.http.get<any>(this.baseUrl.replace('{operation}', 'port').replace('{cluster_name}', clusterName) + "/" + name);
    }
}
