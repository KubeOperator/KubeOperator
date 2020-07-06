import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {ChartMap} from "./chart";

@Injectable({
    providedIn: 'root'
})
export class ChartmuseumService {

    constructor(private http: HttpClient) {
    }

    baseUrl = '/proxy/chartmuseum/{cluster_name}/api/charts/';

    list(clusterName: string): Observable<ChartMap> {
        return this.http.get<ChartMap>(this.baseUrl.replace('{cluster_name}', clusterName));
    }

}
