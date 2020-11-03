import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {BackupAccount, BackupAccountCreateRequest} from './backup-account';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Page} from '../../../shared/class/Page';
import {Host} from '../../host/host';

@Injectable({
    providedIn: 'root'
})
export class BackupAccountService extends BaseModelService<BackupAccount> {
    baseUrl = '/api/v1/backupaccounts';

    constructor(http: HttpClient) {
        super(http);
    }

    listBuckets(item: BackupAccountCreateRequest): Observable<any> {
        const itemUrl = `${this.baseUrl}/buckets`;
        return this.http.post<any>(itemUrl, item);
    }

    listBy(projectName: string): Observable<Page<BackupAccount>> {
        const itemUrl = `${this.baseUrl}?projectName=${projectName}`;
        return this.http.get<Page<BackupAccount>>(itemUrl);
    }
}
