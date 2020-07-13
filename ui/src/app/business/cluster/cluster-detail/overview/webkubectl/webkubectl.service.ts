import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {WebkubectlToken} from "./webkubectl";

@Injectable({
    providedIn: 'root'
})
export class WebkubectlService {

    constructor(private http: HttpClient) {
    }

    baseUrl = '/api/v1/clusters/webkubectl/{cluster_name}/';

    getToken(clusterName: string): Observable<WebkubectlToken> {
        return this.http.get<WebkubectlToken>(this.baseUrl.replace('{cluster_name}', clusterName));
    }
}
