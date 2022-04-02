import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {SourceSearch, SourceCreate, SourceDelete} from './cluster'

@Injectable({
    providedIn: 'root'
})

export class KubernetesService {
    metricUrl = '/api/v1/clusters/kubernetes/search/metric/{cluster_name}';
    searchUrl = '/api/v1/clusters/kubernetes/search';
    createUrl = '/api/v1/clusters/kubernetes/create';
    deleteUrl = '/api/v1/clusters/kubernetes/delete';
    
    limit = 10;
    continueTokenKey = 'continue';

    constructor(private client: HttpClient) {
    }

    getMetrics(clusterName: string): Observable<any> {
        return this.client.post<any>(this.metricUrl.replace('{cluster_name}', clusterName), {});
    }

    listResource(data: SourceSearch): Observable<any> {
        return this.client.post<any>(this.searchUrl, data);
    }

    createResourceNs(data: SourceCreate): Observable<any> {
        return this.client.post<any>(this.createUrl + "/ns", data);
    }
    createResourceSc(data: SourceCreate): Observable<any> {
        return this.client.post<any>(this.createUrl + "/sc", data);
    }
    createResourcePv(data: SourceCreate): Observable<any> {
        return this.client.post<any>(this.createUrl + "/pv", data);
    }
    createResourceSecret(data: SourceCreate): Observable<any> {
        return this.client.post<any>(this.createUrl + "/secret", data);
    }

    deleteResource(data: SourceDelete): Observable<any> {
        return this.client.post<any>(this.deleteUrl, data);
    }
}
