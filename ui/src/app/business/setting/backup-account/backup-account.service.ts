import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {BackupAccount} from './backup-account';
import {HttpClient} from '@angular/common/http';

@Injectable({
    providedIn: 'root'
})
export class BackupAccountService extends BaseModelService<BackupAccount> {
    baseUrl = '/api/v1/backupAccounts';

    constructor(http: HttpClient) {
        super(http);
    }
}
