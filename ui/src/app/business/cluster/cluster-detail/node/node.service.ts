import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {Node, NodeCreateRequest} from "./node";

@Injectable({
    providedIn: 'root'
})
export class NodeService {

    constructor(private http: HttpClient) {
    }

    baseUrl = '/api/v1/clusters/node/{clusterName}/';

    list(clusterName: string): Observable<Node[]> {
        return this.http.get<Node[]>(this.baseUrl.replace('{clusterName}', clusterName));
    }

    create(clusterName: string, item: NodeCreateRequest): Observable<Node[]> {
        return this.http.post<Node[]>(this.baseUrl.replace('{clusterName}', clusterName), item);
    }
}
