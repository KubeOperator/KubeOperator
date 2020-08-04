import {Component, OnInit} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {BackupAccount} from '../backup-account';
import {BackupAccountService} from '../backup-account.service';

@Component({
    selector: 'app-backup-account-list',
    templateUrl: './backup-account-list.component.html',
    styleUrls: ['./backup-account-list.component.css']
})
export class BackupAccountListComponent extends BaseModelComponent<BackupAccount> implements OnInit {

    constructor(private backupAccountService: BackupAccountService) {
        super(backupAccountService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

}
