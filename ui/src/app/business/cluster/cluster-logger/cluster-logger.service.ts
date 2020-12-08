import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {Log} from "./logger";

@Injectable({
    providedIn: 'root'
})
export class ClusterLoggerService {

    constructor(private http: HttpClient) {
    }

    private clusterLoggerUrl = '/api/v1/clusters/logger/{cluster_name}';
    private clusterNodeLoggerUrl = '/api/v1/clusters/node/logger/{cluster_name}/{node_name}';

    getClusterLog(clusterName: string): Observable<Log> {
        return this.http.get<Log>(this.clusterLoggerUrl.replace('{cluster_name}', clusterName));
    }

    getClusterNodeLog(clusterName: string, nodeName: string): Observable<Log> {
        return this.http.get<Log>(this.clusterNodeLoggerUrl.replace('{cluster_name}', clusterName).replace('{node_name}', nodeName));
    }

    openLogger(clusterName: string, nodeName?: string) {
        window.open(`/ui/logger?clusterName=${clusterName}`, '_blank', 'height=865, width=800, top=0, left=0, toolbar=no, menubar=no, scrollbars=no, resizable=yes,location=no, status=no');
    }
}
