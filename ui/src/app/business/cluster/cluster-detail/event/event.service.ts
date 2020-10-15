import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class EventService {

    baseUrl = '/api/v1/events';


    constructor(private http: HttpClient) {

    }

    changeNpd(clusterName: string, operation: string): Observable<any> {
        const itemUrl = `${this.baseUrl}/npd/${operation}/` + clusterName + '/';
        return this.http.post<any>(itemUrl, {});
    }
}
