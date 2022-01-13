import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, Subject} from 'rxjs';
import {License} from './license';

@Injectable({
    providedIn: 'root'
})
export class LicenseService {

    constructor(private http: HttpClient) {
    }

    baseUrl = '/api/v1/license';
    baseHwUrl = '/api/v1/license/hw';

    licenseQueue = new Subject<License>();
    $licenseQueue = this.licenseQueue.asObservable();

    get(): Observable<License> {
        return this.http.get<License>(this.baseUrl);
    }
    gethw(): Observable<License> {
        return this.http.get<License>(this.baseHwUrl);
    }


    setLicense() {
        return this.http.get<License>(this.baseUrl).subscribe(data => {
            this.licenseQueue.next(data);
        });
    }
}
