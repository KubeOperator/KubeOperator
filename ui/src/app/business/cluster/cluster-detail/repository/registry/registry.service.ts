import {Injectable} from '@angular/core';
import {Observable} from "rxjs";
import {ChartMap} from "../chartmuseum/chart";
import {HttpClient} from "@angular/common/http";
import {Registry, RegistryList} from "./registry";

@Injectable({
    providedIn: 'root'
})
export class RegistryService {

    constructor(private http: HttpClient) {
    }

    listUrl = '/proxy/registry/{cluster_name}/v2/_catalog';
    tagsUrl = '/proxy/registry/{cluster_name}/v2/{registry_name}/tags/list/';

    list(clusterName: string): Observable<RegistryList> {
        return this.http.get<RegistryList>(this.listUrl.replace('{cluster_name}', clusterName));
    }

    listTags(clusterName: string, tagName: string): Observable<Registry> {
        return this.http.get<Registry>(this.tagsUrl.replace('{cluster_name}', clusterName).replace('{registry_name}', tagName));
    }

}
