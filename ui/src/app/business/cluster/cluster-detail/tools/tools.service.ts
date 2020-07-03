import {Injectable} from '@angular/core';
import {Observable} from "rxjs";
import {ClusterTool} from "./tools";
import {HttpClient} from "@angular/common/http";

@Injectable({
    providedIn: 'root'
})
export class ToolsService {

    constructor(private http: HttpClient) {
    }

    baseUrl = '/api/v1/clusters/tool/{cluster_name}/';

    list(clusterName: string): Observable<ClusterTool[]> {
        return this.http.get<ClusterTool[]>(this.baseUrl.replace('{cluster_name}', clusterName));
    }

}
