import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import format from 'date-fns/format';


@Injectable({
    providedIn: 'root'
})
export class LoggingService {

    baseUrl = '/api/v1/proxy/logging/{cluster_name}/{index_name}/_search';

    constructor(private http: HttpClient) {
    }

    search(clusterName: string, namespace: string, container: string, pod: string): Observable<any> {
        const date = format(new Date(), 'yyyy.MM.dd');
        const index = 'logstash-' + date;
        const query = {
            query: {
                bool: {
                    must: [
                        {
                            match: {
                                'kubernetes.pod_name': {
                                    query: pod,
                                },
                            },
                        },
                        {
                            match: {
                                'kubernetes.namespace_name': {
                                    query: namespace,
                                },
                            },
                        },
                        {
                            match: {
                                'kubernetes.container_name': {
                                    query: container,
                                },
                            },
                        }
                    ],
                },
            },
            sort: [
                {'@timestamp': 'desc'},
            ]
        };
        return this.http.post<any>(this.baseUrl.replace('{cluster_name}', clusterName).replace('{index_name}', index), query);
    }
}
