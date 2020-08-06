import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BackupFile, BackupStrategy} from './cluster-backup';
import {Observable} from 'rxjs';
import {Page} from '../../../../shared/class/Page';

@Injectable({
    providedIn: 'root'
})
export class BackupService {

    baseUrl = '/api/v1/cluster/backup';
    fileUrl = '/api/v1/cluster/backup/file';

    constructor(private http: HttpClient) {
    }

    getBy(clusterName: string): Observable<BackupStrategy> {
        const itemUrl = `${this.baseUrl}/strategy/${clusterName}/`;
        return this.http.get<BackupStrategy>(itemUrl);
    }

    submit(item: BackupStrategy): Observable<any> {
        const itemUrl = `${this.baseUrl}/strategy/${item.clusterName}/`;
        return this.http.post<any>(itemUrl, item);
    }

    pageBy(page, size, clusterName: string): Observable<Page<BackupFile>> {
        const itemUrl = `${this.baseUrl}/file/${clusterName}/?pageNum=${page}&pageSize=${size}`;
        return this.http.get<Page<BackupFile>>(itemUrl);
    }
}
