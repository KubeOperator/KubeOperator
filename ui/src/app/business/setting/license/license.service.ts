import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {License} from './license';

@Injectable({
    providedIn: 'root'
})
export class LicenseService {

    constructor(private http: HttpClient) {
    }

    baseUrl = '/api/v1/license';

    get(): Observable<License> {
        return this.http.get<License>(this.baseUrl);
    }
}
