import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';
import {HttpClient} from '@angular/common/http';
import {Istio, IstioHelper} from './istios';

@Injectable({
    providedIn: 'root'
})
export class IstioService {

    constructor(private http: HttpClient) {
    }

    baseUrl = '/api/v1/clusters/istio/{operation}/{cluster_name}';

    list(clusterName: string): Observable<IstioHelper[]> {
        return this.http.get<IstioHelper[]>(this.baseUrl.replace('/{operation}', '').replace('{cluster_name}', clusterName));
    }

    enable(clusterName: string, items: IstioHelper[]): Observable<IstioHelper[]> {
        return this.http.post<IstioHelper[]>(this.baseUrl.replace('{operation}', 'enable').replace('{cluster_name}', clusterName), items);
    }

    disable(clusterName: string, items: IstioHelper[]): Observable<IstioHelper[]> {
        return this.http.post<IstioHelper[]>(this.baseUrl.replace('{operation}', 'disable').replace('{cluster_name}', clusterName), items);
    }
}
