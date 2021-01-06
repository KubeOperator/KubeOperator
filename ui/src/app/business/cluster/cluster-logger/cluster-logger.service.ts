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
    private storageProvisionerLoggerUrl = '/api/v1/clusters/provisioner/log/{cluster_name}/{log_id}';

    getClusterLog(clusterName: string): Observable<Log> {
        return this.http.get<Log>(this.clusterLoggerUrl.replace('{cluster_name}', clusterName));
    }

    getClusterNodeLog(clusterName: string, nodeName: string): Observable<Log> {
        return this.http.get<Log>(this.clusterNodeLoggerUrl.replace('{cluster_name}', clusterName).replace('{node_name}', nodeName));
    }

    getProvisionerLog(clusterName: string, logId: string): Observable<Log> {
        return this.http.get<Log>(this.storageProvisionerLoggerUrl.replace('{cluster_name}', clusterName).replace('{log_id}', logId));
    }

    openLogger(clusterName: string, nodeName?: string) {
        window.open(`/ui/logger?clusterName=${clusterName}`, '_blank', 'height=865, width=800, top=0, left=0, toolbar=no, menubar=no, scrollbars=no, resizable=yes,location=no, status=no');
    }

    openProvisionerLogger(clusterName: string, logId?: string) {
        window.open(`/ui/logger?clusterName=${clusterName}&logId=${logId}`, '_blank', 'height=865, width=800, top=0, left=0, toolbar=no, menubar=no, scrollbars=no, resizable=yes,location=no, status=no');
    }
}
