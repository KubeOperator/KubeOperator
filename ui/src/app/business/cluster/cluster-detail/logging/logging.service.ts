import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';


@Injectable({
    providedIn: 'root'
})
export class LoggingService {

    efBaseUrl = '/proxy/logging/{cluster_name}/{index_name}/_search?pretty=true';
    lokiBaseUrl = '/proxy/loki/{cluster_name}/';
    constructor(private http: HttpClient) {
    }
    EfSearch(clusterName: string, queryArry: any[], queryIndex: string,
           beginDate: string, endDate: string, pageFrom: number, pageSize: number): Observable<any> {
        const index = queryIndex;
        const query = {
            from: (pageFrom - 1) * pageSize,
            size: pageSize,
            query: {
                bool: {
                    must: queryArry,
                    filter: {
                        range: {
                            '@timestamp': {
                                gte: beginDate,
                                lte: endDate,
                                format: 'yyyy.MM.dd',
                                time_zone: '+08:00'
                            }
                        }
                    }
                }
            },
            sort: [
                {'@timestamp': 'desc'},
            ]
        };
        return this.http.post<any>(this.efBaseUrl.replace('{cluster_name}', clusterName).replace('{index_name}', index) + '&ignore_unavailable=true', query);
    }
    LokiLabels(clusterName: string): Observable<any> {
        return this.http.post<any>(this.lokiBaseUrl.replace('{cluster_name}', clusterName) + 'loki/api/v1/labels', '');
    }
    LokiLabelValues(clusterName: string, label: string): Observable<any> {
        return this.http.post<any>(this.lokiBaseUrl.replace('{cluster_name}', clusterName) + 'loki/api/v1/label/{label}/values'.replace('{label}', label), '');
    }
    LokiSearch(clusterName: string, params: string): Observable<any> {
        return this.http.post<any>(this.lokiBaseUrl.replace('{cluster_name}', clusterName) + 'loki/api/v1/query_range?' + params, '');
    }
}
