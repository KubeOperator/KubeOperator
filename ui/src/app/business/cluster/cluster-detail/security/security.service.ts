import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {CisTask} from "./security";
import {Page} from "../../../../shared/class/Page";

@Injectable({
    providedIn: 'root'
})
export class SecurityService {

    constructor(private http: HttpClient) {
    }

    private url = '/api/v1/clusters/cis/{cluster_name}';

    page(clusterName: string, page: number, size: number): Observable<Page<CisTask>> {
        return this.http.get<Page<CisTask>>(this.url.replace('{cluster_name}', clusterName) + `?pageNum=${page}&pageSize=${size}`);
    }

    create(clusterName: string): Observable<CisTask> {
        return this.http.post<CisTask>(this.url.replace('{cluster_name}', clusterName), {});
    }

    delete(clusterName: string, id: string): Observable<any> {
        return this.http.delete(this.url.replace('{cluster_name}', clusterName) + `/${id}`);
    }
}
