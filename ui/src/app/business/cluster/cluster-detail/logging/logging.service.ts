import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';


@Injectable({
    providedIn: 'root'
})
export class LoggingService {

    baseUrl = '/proxy/logging/{cluster_name}/{index_name}/_search?pretty=true';
    constructor(private http: HttpClient) {
    }
    Search(clusterName: string, queryArry: any[], queryIndex: string,
           beginDate: string, endDate: string, pageFrom: number, pageSize: number): Observable<any> {
        const index = queryIndex;
        const query = {
            from: pageFrom,
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
        return this.http.post<any>(this.baseUrl.replace('{cluster_name}', clusterName).replace('{index_name}', index) + '&ignore_unavailable=true', query);
    }
}
