import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../../shared/class/BaseModelService';
import {HttpClient} from '@angular/common/http';
import {BackupFile, BackupStrategy} from './cluster-backup';
import {Observable} from 'rxjs';
import {Page} from '../../../../shared/class/Page';

@Injectable({
    providedIn: 'root'
})
export class BackupFileService extends BaseModelService<BackupFile> {

    baseUrl = '/api/v1/clusters/backup/files';

    constructor(http: HttpClient) {
        super(http);
    }

    pageBy(page, size, clusterName: string): Observable<Page<BackupFile>> {
        const itemUrl = `${this.baseUrl}?pageNum=${page}&pageSize=${size}&clusterName=${clusterName}`;
        return this.http.get<Page<BackupFile>>(itemUrl);
    }


    backup(item: BackupFile): Observable<any> {
        const itemUrl = `${this.baseUrl}/backup`;
        return this.http.post<any>(itemUrl, item);
    }

    restore(item: BackupFile): Observable<any> {
        const itemUrl = `${this.baseUrl}/restore`;
        return this.http.post<any>(itemUrl, item);
    }

    localRestore(formData): Observable<any> {
        const itemUrl = `${this.baseUrl}/restore/local`;
        return this.http.post<any>(itemUrl, formData);
    }

}
