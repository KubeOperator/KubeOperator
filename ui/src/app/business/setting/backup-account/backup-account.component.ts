import {Component, OnInit, ViewChild} from '@angular/core';
import {BackupAccountListComponent} from './backup-account-list/backup-account-list.component';
import {BackupAccountCreateComponent} from './backup-account-create/backup-account-create.component';
import {BackupAccountUpdateComponent} from './backup-account-update/backup-account-update.component';
import {BackupAccountDeleteComponent} from './backup-account-delete/backup-account-delete.component';

@Component({
    selector: 'app-backup-account',
    templateUrl: './backup-account.component.html',
    styleUrls: ['./backup-account.component.css']
})
export class BackupAccountComponent implements OnInit {

    @ViewChild(BackupAccountListComponent, {static: true})
    list: BackupAccountListComponent;

    @ViewChild(BackupAccountCreateComponent, {static: true})
    create: BackupAccountCreateComponent;

    @ViewChild(BackupAccountUpdateComponent, {static: true})
    update: BackupAccountUpdateComponent;

    @ViewChild(BackupAccountDeleteComponent, {static: true})
    delete: BackupAccountDeleteComponent;

    constructor() {
    }

    ngOnInit(): void {

    }

    refresh() {
        this.list.reset();
        this.list.refresh();
    }

    onCreate() {
        this.create.open();
    }

    onDelete(items) {
        this.delete.open(items);
    }

    onUpdate(item) {
        this.update.open(item);
    }
}
