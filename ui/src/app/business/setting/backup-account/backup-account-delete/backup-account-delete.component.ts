import {Component, OnInit} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {BackupAccount} from '../backup-account';
import {BackupAccountService} from '../backup-account.service';

@Component({
    selector: 'app-backup-account-delete',
    templateUrl: './backup-account-delete.component.html',
    styleUrls: ['./backup-account-delete.component.css']
})
export class BackupAccountDeleteComponent extends BaseModelComponent<BackupAccount> implements OnInit {
    opened = false;
    items: BackupAccount[] = [];

    constructor(private backupAccountService: BackupAccountService) {
        super(backupAccountService);
    }

    ngOnInit(): void {
    }

    open(items) {
        this.items = items;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
        this.items = [];
    }

    onSubmit() {

    }
}
