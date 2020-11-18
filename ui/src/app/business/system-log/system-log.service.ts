import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {SystemLog, Page} from "./system-log";
import {Observable} from "rxjs";

@Injectable({
    providedIn: 'root'
})
export class SystemLogService {
    constructor(private http: HttpClient) {
    }
    baseUrl = '/api/v1/logs';

    list(page, size): Observable<Page<SystemLog>> {
        const url = this.baseUrl
        const pageUrl = `${url}?pageNum=${page}&pageSize=${size}`;
        return this.http.get<Page<SystemLog>>(pageUrl);
    }
}
