import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {BackupAccountComponent} from './backup-account.component';
import {BackupAccountCreateComponent} from './backup-account-create/backup-account-create.component';
import {BackupAccountListComponent} from './backup-account-list/backup-account-list.component';
import {BackupAccountUpdateComponent} from './backup-account-update/backup-account-update.component';
import {BackupAccountDeleteComponent} from './backup-account-delete/backup-account-delete.component';
import {CoreModule} from '../../../core/core.module';


@NgModule({
    declarations: [BackupAccountComponent, BackupAccountCreateComponent, BackupAccountListComponent,
        BackupAccountUpdateComponent, BackupAccountDeleteComponent],
    imports: [
        CommonModule,
        CoreModule
    ]
})
export class BackupAccountModule {
}
