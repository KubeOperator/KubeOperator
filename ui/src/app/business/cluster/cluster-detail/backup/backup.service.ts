import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BackupStrategy} from './cluster-backup';
import {Observable} from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class BackupService {

    baseUrl = '/api/v1/clusters/backup/strategy';

    constructor(private http: HttpClient) {
    }

    getBy(clusterName: string): Observable<BackupStrategy> {
        const itemUrl = `${this.baseUrl}/${clusterName}`;
        return this.http.get<BackupStrategy>(itemUrl);
    }

    submit(item: BackupStrategy): Observable<any> {
        const itemUrl = `${this.baseUrl}`;
        return this.http.post<any>(itemUrl, item);
    }
}
