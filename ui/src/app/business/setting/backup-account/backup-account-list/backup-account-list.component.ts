import {Component, OnInit} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {BackupAccount} from '../backup-account';
import {BackupAccountService} from '../backup-account.service';

@Component({
    selector: 'app-backup-account-list',
    templateUrl: './backup-account-list.component.html',
    styleUrls: ['./backup-account-list.component.css']
})
export class BackupAccountListComponent extends BaseModelDirective<BackupAccount> implements OnInit {

    constructor(private backupAccountService: BackupAccountService) {
        super(backupAccountService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

}
