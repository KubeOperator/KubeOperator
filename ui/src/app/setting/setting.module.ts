import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {SettingComponent} from './setting.component';

import {CoreModule} from '../core/core.module';
import {SystemSettingComponent} from './system-setting/system-setting.component';
import {BackupStorageSettingComponent } from './backup-storage-setting/backup-storage-setting.component';
import {ClusterBackupModule} from '../cluster-backup/cluster-backup.module';
import {BackupStorageSettingModule} from './backup-storage-setting/backup-storage-setting.module';

@NgModule({
  declarations: [SettingComponent, SystemSettingComponent, BackupStorageSettingComponent],
  imports: [
    CommonModule,
    CoreModule,

    ClusterBackupModule,
    BackupStorageSettingModule,

  ]
})
export class SettingModule {
}
