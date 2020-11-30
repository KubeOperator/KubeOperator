import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {BackupAccountComponent} from './backup-account.component';
import {BackupAccountCreateComponent} from './backup-account-create/backup-account-create.component';
import {BackupAccountListComponent} from './backup-account-list/backup-account-list.component';
import {BackupAccountUpdateComponent} from './backup-account-update/backup-account-update.component';
import {BackupAccountDeleteComponent} from './backup-account-delete/backup-account-delete.component';
import {BackupAccountGrantComponent} from './backup-account-grant/backup-account-grant.component';
import {CoreModule} from '../../../core/core.module';
import {SharedModule} from '../../../shared/shared.module';


@NgModule({
    declarations: [BackupAccountComponent, BackupAccountCreateComponent, BackupAccountListComponent,
        BackupAccountUpdateComponent, BackupAccountDeleteComponent, BackupAccountGrantComponent],
    imports: [
        CommonModule,
        CoreModule,
        SharedModule
    ]
})
export class BackupAccountModule {
}
