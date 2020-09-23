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

    private url = '/api/v1/clusters/logger/{cluster_name}';

    get(clusterName: string): Observable<Log> {
        return this.http.get<Log>(this.url.replace('{cluster_name}', clusterName));
    }
    openLogger(clusterName: string) {
        window.open('/ui/logger?clusterName=' + clusterName, 'blank', 'height=820, width=800, top=0, left=0, toolbar=no, menubar=no, scrollbars=no, resizable=yes,location=no, status=no');
    }
}
